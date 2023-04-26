// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

type AddMemoryInput struct {
	// 标题
	Title string `json:"title"`
	// 内容
	Content string `json:"content"`
}

type Comment struct {
	// Comment ID
	ID string `json:"id"`
	// Memory ID
	MemoryID string `json:"memory_id"`
	// 创建用户
	User *User `json:"user"`
	// 内容
	Content string `json:"content"`
	// 是否已归档
	Archived bool `json:"archived"`
	// 发布时间
	CreateTime int64 `json:"create_time"`
	// 修改时间
	UpdateTime int64 `json:"update_time"`
}

type CommentInput struct {
	// 内容
	Content string `json:"content"`
}

type ForgetInput struct {
	Email string `json:"email"`
	// 新密码
	Password        string `json:"password"`
	EmailVerifyCode string `json:"email_verify_code"`
}

type HashTag struct {
	// HashTag ID
	ID string `json:"id"`
	// Memory ID
	MemoryID string `json:"memory_id"`
	// 创建用户
	User *User `json:"user"`
	// Tag名称
	Name string `json:"name"`
	// 是否已归档
	Archived bool `json:"archived"`
	// 创建时间
	CreateTime int64 `json:"create_time"`
	// 修改时间
	UpdateTime int64 `json:"update_time"`
}

type HashTagInput struct {
	// 名称
	Name string `json:"name"`
}

type History struct {
	// History ID
	ID string `json:"id"`
	// Memory ID
	MemoryID string `json:"memory_id"`
	// 创建用户
	User *User `json:"user"`
	// 标题
	Title string `json:"title"`
	// 内容
	Content string `json:"content"`
	// 发布时间
	CreateTime int64 `json:"create_time"`
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

type Memory struct {
	// Memory ID
	ID string `json:"id"`
	// 创建用户
	User *User `json:"user"`
	// 标题
	Title string `json:"title"`
	// 内容
	Content  string     `json:"content"`
	Hashtags []*HashTag `json:"hashtags"`
	// 是否已归档
	Archived bool `json:"archived"`
	// 发布时间
	CreateTime int64 `json:"create_time"`
	// 修改时间
	UpdateTime int64 `json:"update_time"`
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
	User *User `json:"user"`
	// 大小(Byte)
	Size int `json:"size"`
	// 引用该资源的 Memories
	Memories []*Memory `json:"memories"`
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

type UploadTokenPayload struct {
	// 上传凭证
	Token string `json:"token"`
	// 相对路径
	Path string `json:"path"`
	// URL
	URL string `json:"url"`
}

type User struct {
	// 用户ID
	ID string `json:"id"`
	// 用户名
	Username string `json:"username"`
	// 头像URL
	Avatar string `json:"avatar"`
	// 邮箱
	Mail string `json:"mail"`
	// 上次登录时间
	LoginTime int64 `json:"login_time"`
	// 注册时间
	CreateTime int64 `json:"create_time"`
	// 权限
	Permission int `json:"permission"`
	// 已使用资源(Byte)
	Used int `json:"used"`
	// 订阅
	Subscribe *Subscribe `json:"subscribe"`
}
