package storage

import (
	"github.com/go-redis/redis/v8"
	"memox_server/src/service/storage/bce"
	"memox_server/src/service/storage/local"
	"memox_server/src/service/storage/utils"
	"strings"
)

func NewStorageSvc(conf Config, r *redis.Client) utils.Service {
	var service utils.Service
	provider := strings.ToLower(conf.StorageProvider)
	switch provider {
	case "bce":
		service = bce.NewBCESvc(conf.BCE)
	case "local":
		service = local.NewLocalSvc(conf.Local, r)
	default:
		panic("unknown storage provider")
	}
	return service
}
