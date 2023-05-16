package utils

import (
	"context"
)

type Service interface {
	// GetToken 获取上传凭证
	GetToken(ctx context.Context, fileName string) (*UploadTokenPayload, error)
	// GetUrl 获取文件链接
	GetUrl(ctx context.Context, fileName string) (string, error)
	// Delete 删除文件
	Delete(ctx context.Context, fileName string) (bool, error)
}

type UploadTokenPayload struct {
	// 唯一资源标识符
	ID string `json:"id"`
	// 是否已存在
	Exist bool `json:"exist"`
	// 用于STS凭证访问的AK
	AccessKey string `json:"access_key"`
	// 用于STS凭证访问的SK
	SecretAccessKey string `json:"secret_access_key"`
	// SessionToken，使用STS凭证访问时必须携带
	SessionToken string `json:"session_token"`
	// UserId
	UserID string `json:"user_id"`
	// 文件名
	FileName string `json:"file_name"`
}
