package subscribe

type Config struct {
	DefaultCapacity      int64  `yaml:"default_capacity"`       // 默认容量 (Byte)
	DefaultSubscribeName string `yaml:"default_subscribe_name"` // 默认订阅名称
}
