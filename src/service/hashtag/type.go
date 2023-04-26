package hashtag

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hashtag struct {
	ObjectID   primitive.ObjectID `bson:"_id"`         // ID
	Uid        primitive.ObjectID `bson:"uid"`         // 用户UID
	Name       string             `bson:"name"`        // 名称
	Archived   bool               `bson:"archived"`    // 归档，归档后话题列表中不可见，但仍可通过Memories的引用查看
	CreateTime int64              `bson:"create_time"` // 创建时间
}
