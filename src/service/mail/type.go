package mail

// Code 验证码
type Code struct {
	ID   string `json:"id"`   // 验证码标志
	Code string `json:"code"` // 验证码
	Act  string `json:"act"`  // 行为：注册/找回
}
