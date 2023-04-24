package mail

import "errors"

type Config struct {
	SmtpMailHost     string `yaml:"smtp_mail_host"`
	SmtpMailPort     int    `yaml:"smtp_mail_port"`
	SmtpMailUser     string `yaml:"smtp_mail_user"`
	SmtpMailPwd      string `yaml:"smtp_mail_pwd"`
	SmtpMailNickname string `yaml:"smtp_mail_nickname"`
	Subject          string `yaml:"subject"`
	Template         string `yaml:"template"`
	CodeExpire       int    `yaml:"code_expire"`
	CodeLength       int    `yaml:"code_length"`
	CodeCoolDown     int    `yaml:"code_cool_down"`
}

const (
	keyVerifyCode         = "verify_code:%s"
	keyVerifyCodeCoolDown = "verify_code_cd:%s"
)

var ErrVerifyCodeCoolDown = errors.New("request verify code id is in cool down time")
