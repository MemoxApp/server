package memory

import "go.mongodb.org/mongo-driver/bson/primitive"

type Memory struct {
	ObjectID   primitive.ObjectID   `bson:"_id"`         // ID
	Uid        primitive.ObjectID   `bson:"uid"`         // 用户UID
	Title      string               `bson:"title"`       // 标题
	Content    string               `bson:"content"`     // 内容
	HashTags   []primitive.ObjectID `bson:"hash_tags"`   // 话题
	Archived   bool                 `bson:"archived"`    // 归档 true 已归档 false 正常
	CreateTime int64                `bson:"create_time"` // 创建时间
	UpdateTime int64                `bson:"update_time"` // 修改时间
}
