package bce

import (
	"context"
	"fmt"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/sts"
	"time_speak_server/src/exception"
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
	exist, err := bosClient.DoesBucketExist(config.BucketName)
	if err != nil {
		panic(err)
	}
	if !exist {
		// 创建Bucket
		location, err := bosClient.PutBucket(config.BucketName)
		if err != nil {
			fmt.Println("创建存储桶失败", err)
			panic(err)
		}
		fmt.Println("成功创建存储桶，路径：", location)
	}
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
	if !utils.CheckFileName(fileName) {
		return nil, exception.ErrInvalidFileName
	}
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	// 生成上传凭证
	aclString := b.getWritePermissionACL(userId.Hex(), fileName)
	sessionToken, err := b.Sts.GetSessionToken(60, aclString)
	if err != nil {
		return nil, err
	}
	return &utils.UploadTokenPayload{
		AccessKey:       sessionToken.AccessKeyId,
		SecretAccessKey: sessionToken.SecretAccessKey,
		SessionToken:    sessionToken.SessionToken,
		UserID:          sessionToken.UserId,
	}, nil
}

func (b *BCE) GetUrl(ctx context.Context, path string) (string, error) {
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return "", err
	}
	p := utils.GeneratePath(userId.Hex(), path)
	// 生成下载地址,有效时间 30 min
	url := b.Bos.BasicGeneratePresignedUrl(b.Config.BucketName, p, 1800)
	return url, nil
}

func (b *BCE) Delete(ctx context.Context, path string) (bool, error) {
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return false, err
	}
	p := utils.GeneratePath(userId.Hex(), path)
	// 删除文件
	err = b.Bos.DeleteObject(b.Config.BucketName, p)
	if err != nil {
		return false, err
	}
	return true, nil
}
