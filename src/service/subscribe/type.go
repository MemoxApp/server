package subscribe

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subscribe struct {
	ObjectID   primitive.ObjectID `bson:"_id"`         // ID
	Name       string             `bson:"name"`        // 名称
	Capacity   int64              `bson:"capacity"`    // 容量 (Byte)
	Enabled    bool               `bson:"enabled"`     // 启用 true 已归档 false 正常
	CreateTime int64              `bson:"create_time"` // 创建时间
	UpdateTime int64              `bson:"update_time"` // 修改时间
}

const KB = 1024
const MB = 1024 * KB
const GB = 1024 * MB
