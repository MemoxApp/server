package bce

import (
	"context"
	"errors"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/sts"
	"memox_server/src/log"
	"memox_server/src/service/storage/utils"
	"memox_server/src/service/user"
	"strings"
	"time"
)

type BCE struct {
	Config Config
}

// NewBCESvc 创建BCE服务
func NewBCESvc(config Config) *BCE {
	// 创建BOS服务的Client
	return &BCE{
		Config: config,
	}
}

func (b *BCE) GetToken(ctx context.Context, fileName string) (*utils.UploadTokenPayload, error) {
	userId, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	// 生成上传凭证
	aclString := b.getWritePermissionACL(userId.Hex(), fileName, b.Config.BucketName)
	absPath := utils.GenerateResourcePath(userId.Hex(), fileName)
	client, err := sts.NewClient(b.Config.AccessKeyID, b.Config.SecretAccessKey)
	sessionToken, err := client.GetSessionToken(60, aclString)
	if err != nil {
		log.Error("GetSessionToken Error, Time:" + time.Now().String() + ", error:" + err.Error())
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
		if strings.ToLower(b.Config.CdnAuthType) == "a" {
			u := SignUrlTypeA(b.Config.EndPoint, p, userId.Hex(), b.Config.CdnAuthKey)
			return u, nil
		} else if strings.ToLower(b.Config.CdnAuthType) == "b" {
			u := SignUrlTypeB(b.Config.EndPoint, p, b.Config.CdnAuthKey)
			return u, nil
		} else if strings.ToLower(b.Config.CdnAuthType) == "c" {
			u := SignUrlTypeC(b.Config.EndPoint, p, b.Config.CdnAuthKey)
			return u, nil
		}
		return "", errors.New("no specified cdn auth type")
	} else {
		// 生成对象存储下载地址,有效时间 30 min
		bosClient, err := bos.NewClient(b.Config.AccessKeyID, b.Config.SecretAccessKey, b.Config.EndPoint)
		if err != nil {
			return "", err
		}
		u := bosClient.BasicGeneratePresignedUrl(b.Config.BucketName, p, 1800)
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
	bosClient, err := bos.NewClient(b.Config.AccessKeyID, b.Config.SecretAccessKey, b.Config.EndPoint)
	if err != nil {
		return false, err
	}
	err = bosClient.DeleteObject(b.Config.BucketName, p)
	if err != nil {
		return false, err
	}
	return true, nil
}
