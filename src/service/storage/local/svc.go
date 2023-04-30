package local

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-redis/redis/v8"
	"io"
	"os"
	"path/filepath"
	"time"
	"time_speak_server/src/exception"
	"time_speak_server/src/service/storage/utils"
	"time_speak_server/src/service/user"
)

type Local struct {
	Config Config
	r      *redis.Client
}

// NewLocalSvc 创建本地文件服务
func NewLocalSvc(config Config, r *redis.Client) *Local {
	err := os.MkdirAll(config.Folder, 0666)
	if err != nil {
		panic(err)
	}
	return &Local{
		Config: config,
		r:      r,
	}
}

func (b *Local) GetToken(ctx context.Context, fileName string) (*utils.UploadTokenPayload, error) {
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	randStr := GenerateRandomString(32)
	b.r.Set(ctx, "Local-"+randStr, fileName, time.Minute)
	return &utils.UploadTokenPayload{
		AccessKey:       "",
		SecretAccessKey: "",
		SessionToken:    randStr,
		UserID:          userId.Hex(),
	}, nil
}

func (b *Local) GetUrl(ctx context.Context, path string) (string, error) {
	p := b.Config.Schema + "://" + b.Config.Host + "/resources/" + path
	return p, nil
}

func (b *Local) Delete(ctx context.Context, path string) (bool, error) {
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return false, err
	}
	p := utils.GenerateResourcePath(userId.Hex(), path)
	// 删除文件
	absPath := filepath.Clean(b.Config.Folder + "/" + p)
	_ = os.Remove(absPath)
	return true, nil
}

func (b *Local) Upload(ctx context.Context, session string, upload graphql.Upload) (string, int64, error) {
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return "", 0, err
	}
	s, err := b.r.Get(ctx, "Local-"+session).Result()
	if err != nil {
		return "", 0, exception.ErrInvalidSession
	}
	path := utils.GenerateResourcePath(userId.Hex(), s)
	absPath := filepath.Clean(b.Config.Folder + "/" + path)
	err = os.MkdirAll(filepath.Dir(absPath), 0666)
	if err != nil {
		return "", 0, err
	}
	f, err := os.OpenFile(absPath, os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", 0, err
	}
	// 关闭文件
	defer f.Close()
	// 写入文件
	var buffer = make([]byte, 1024)
	var size int64
	for {
		n, err := upload.File.Read(buffer)
		size += int64(n)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", 0, err
		}
		_, err = f.Write(buffer[:n])
		if err != nil {
			return "", 0, err
		}
	}
	_, _ = b.r.Del(ctx, "Local-"+session).Result()
	return s, size, nil
}
