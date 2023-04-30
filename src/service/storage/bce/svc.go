package bce

import (
	"context"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/sts"
	"time_speak_server/src/service/storage/utils"
	"time_speak_server/src/service/user"
)

type BCE struct {
	Config Config
	Bos    *bos.Client
	Sts    *sts.Client
}

// NewBCESvc 创建BCE服务
func NewBCESvc(config Config) *BCE {
	// 创建BOS服务的Client
	bosClient, err := bos.NewClient(config.AccessKeyID, config.SecretAccessKey, config.EndPoint)
	client, err := sts.NewClient(config.AccessKeyID, config.SecretAccessKey)
	if err != nil {
		return nil
	}
	return &BCE{
		Config: config,
		Bos:    bosClient,
		Sts:    client,
	}
}

func (b *BCE) GetToken(ctx context.Context, fileName string) (*utils.UploadTokenPayload, error) {
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	// 生成上传凭证
	aclString := b.getWritePermissionACL(userId.Hex(), fileName)
	absPath := utils.GenerateResourcePath(userId.Hex(), fileName)
	sessionToken, err := b.Sts.GetSessionToken(60, aclString)
	if err != nil {
		return nil, err
	}
	return &utils.UploadTokenPayload{
		AccessKey:       sessionToken.AccessKeyId,
		SecretAccessKey: sessionToken.SecretAccessKey,
		SessionToken:    sessionToken.SessionToken,
		UserID:          sessionToken.UserId,
		FileName:        absPath,
	}, nil
}

func (b *BCE) GetUrl(ctx context.Context, path string) (string, error) {
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return "", err
	}
	p := utils.GenerateResourcePath(userId.Hex(), path)
	if b.Config.CDN {
		// 生成CDN下载地址,有效时间在CDN控制台配置
		u := SignUrl(b.Config.EndPoint, p, b.Config.CdnAuthKey)
		return u, nil
	} else {
		// 生成对象存储下载地址,有效时间 30 min
		u := b.Bos.BasicGeneratePresignedUrl(b.Config.BucketName, p, 1800)
		return u, nil
	}
}

func (b *BCE) Delete(ctx context.Context, path string) (bool, error) {
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return false, err
	}
	p := utils.GenerateResourcePath(userId.Hex(), path)
	// 删除文件
	err = b.Bos.DeleteObject(b.Config.BucketName, p)
	if err != nil {
		return false, err
	}
	return true, nil
}
