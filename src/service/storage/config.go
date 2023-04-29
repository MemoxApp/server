package storage

import (
	"time_speak_server/src/service/storage/bce"
	"time_speak_server/src/service/storage/local"
)

type Config struct {
	StorageProvider string       `yaml:"storage_provider"`
	BCE             bce.Config   `yaml:"bce"`
	Local           local.Config `yaml:"local"`
}
