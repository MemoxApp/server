package resource

import "go.mongodb.org/mongo-driver/bson/primitive"

type Resource struct {
	ObjectID   primitive.ObjectID   `bson:"_id"`         // ID
	Uid        primitive.ObjectID   `bson:"uid"`         // 用户UID
	Path       string               `bson:"path"`        // 内容
	Size       int64                `bson:"size"`        // 大小
	Ref        []primitive.ObjectID `bson:"ref"`         // 引用资源的 Memory ID
	CreateTime int64                `bson:"create_time"` // 创建时间
	UpdateTime int64                `bson:"update_time"` // 更新时间
}
