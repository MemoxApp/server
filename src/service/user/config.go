package user

type Config struct {
	TokenExpire int    `yaml:"token_expire"`
	TokenSecret string `yaml:"token_secret"`
}
