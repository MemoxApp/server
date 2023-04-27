package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time_speak_server/src/log"
	"time_speak_server/src/service/comment"
	"time_speak_server/src/service/hashtag"
	"time_speak_server/src/service/history"
	"time_speak_server/src/service/mail"
	"time_speak_server/src/service/memory"
	"time_speak_server/src/service/user"
)

type ContextKey string

type Config struct {
	App     app            `yaml:"app"`
	DB      db             `yaml:"db"`
	Debug   bool           `yaml:"debug"`
	User    user.Config    `yaml:"user"`
	Mail    mail.Config    `yaml:"mail"`
	Memory  memory.Config  `yaml:"memory"`
	History history.Config `yaml:"history"`
	Hashtag hashtag.Config `yaml:"hashtag"`
	Comment comment.Config `yaml:"comment"`
}

type app struct {
	Addr   string `yaml:"addr"`
	Host   string `yaml:"host"`
	Prefix string `yaml:"prefix"`
}

type db struct {
	MongoAddr string `yaml:"mongo_addr"`
	MongoDB   string `yaml:"mongo_db"`
	RedisAddr string `yaml:"redis_addr"`
	RedisDB   int    `yaml:"redis_db"`
}

// MustReadConfigFile 从指定文件读取配置
func MustReadConfigFile(filename string) Config {
	data, err := os.ReadFile(filename)

	if err != nil {
		log.Fatal("read config error", "filename", filename, "err", err)
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("unmarshal config error", "filename", filename, "err", err)
	}

	log.Info("successfully read config file", "filename", filename)
	return config
}

// WriteConfigFile 将配置写入指定文件
func WriteConfigFile(c Config, filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		log.Error("marshal config error", "filename", filename, "err", err)
		return err
	}

	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		log.Error("write config error", "filename", filename, "err", err)
	}
	return err
}
