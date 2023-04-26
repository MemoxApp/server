package history

import "go.mongodb.org/mongo-driver/bson/primitive"

type History struct {
	ObjectID   primitive.ObjectID   `bson:"_id"`         // ID
	Uid        primitive.ObjectID   `bson:"uid"`         // 用户UID
	MemoryID   primitive.ObjectID   `bson:"memory_id"`   // Memory ID
	Title      string               `bson:"title"`       // 标题
	Content    string               `bson:"content"`     // 内容
	HashTags   []primitive.ObjectID `bson:"hash_tags"`   // 话题
	CreateTime int64                `bson:"create_time"` // 创建时间
}
