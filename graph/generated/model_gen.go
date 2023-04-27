// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"time_speak_server/src/service/memory"
	"time_speak_server/src/service/user"
)

type AddCommentInput struct {
	// Comment 对象ID
	ID string `json:"id"`
	// 是否子回复
	SubComment bool `json:"subComment"`
	// 内容
	Content string `json:"content"`
}

type AddMemoryInput struct {
	// 标题
	Title string `json:"title"`
	// 内容
	Content string `json:"content"`
}

type ForgetInput struct {
	Email string `json:"email"`
	// 新密码
	Password        string `json:"password"`
	EmailVerifyCode string `json:"email_verify_code"`
}

type HashTagInput struct {
	// ID
	ID string `json:"id"`
	// 名称
	Name *string `json:"name,omitempty"`
	// 是否已归档
	Archived *bool `json:"archived,omitempty"`
}

type ListInput struct {
	Page     int  `json:"page"`
	Size     int  `json:"size"`
	ByCreate bool `json:"byCreate"`
	Desc     bool `json:"desc"`
	Archived bool `json:"archived"`
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginPayload struct {
	ID         string `json:"id"`
	Token      string `json:"token"`
	Permission int    `json:"permission"`
	Expire     int64  `json:"expire"`
}

type RegisterInput struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	EmailVerifyCode string `json:"email_verify_code"`
}

type Resource struct {
	// Resource ID
	ID string `json:"id"`
	// 资源路径
	Path string `json:"path"`
	// 创建用户
	User *user.User `json:"user"`
	// 大小(Byte)
	Size int `json:"size"`
	// 引用该资源的 Memories
	Memories []*memory.Memory `json:"memories"`
	// 是否已归档
	Archived bool `json:"archived"`
	// 创建时间
	CreateTime int64 `json:"create_time"`
}

type ResourceInput struct {
	// 文件名
	FileName string `json:"file_name"`
}

type SendEmailCodeInput struct {
	Mail     string `json:"mail"`
	Register bool   `json:"register"`
}

type Subscribe struct {
	// Subscribe ID
	ID string `json:"id"`
	// 订阅名称
	Name string `json:"name"`
	// 资源额度(Byte)
	Capacity int `json:"capacity"`
	// 是否启用
	Available string `json:"available"`
	// 创建时间
	CreateTime int64 `json:"create_time"`
	// 修改时间
	UpdateTime int64 `json:"update_time"`
}

type SubscribeInput struct {
	// 订阅名称
	Name string `json:"name"`
	// 资源额度(Byte)
	Capacity int `json:"capacity"`
	// 是否启用
	Available string `json:"available"`
}

type UpdateCommentInput struct {
	// Comment ID
	ID string `json:"id"`
	// 内容
	Content *string `json:"content,omitempty"`
	// 是否已归档
	Archived *bool `json:"archived,omitempty"`
}

type UpdateMemoryInput struct {
	// ID
	ID string `json:"id"`
	// 标题
	Title string `json:"title"`
	// 内容
	Content string `json:"content"`
}

type UploadTokenPayload struct {
	// 上传凭证
	Token string `json:"token"`
	// 相对路径
	Path string `json:"path"`
	// URL
	URL string `json:"url"`
}
