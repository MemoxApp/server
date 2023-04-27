package comment

import "go.mongodb.org/mongo-driver/bson/primitive"

type Comment struct {
	ObjectID   primitive.ObjectID   `bson:"_id"`         // ID
	Uid        primitive.ObjectID   `bson:"uid"`         // 用户UID
	CommentID  primitive.ObjectID   `bson:"comment_id"`  // 评论对象 ID (Top-level 回复为 Memory ID)
	ParentID   primitive.ObjectID   `bson:"parent_id"`   // 父评论 ID (Top-level 回复为 Memory ID)
	Content    string               `bson:"content"`     // 内容
	HashTags   []primitive.ObjectID `bson:"hash_tags"`   // 话题
	Archived   bool                 `bson:"archived"`    // 归档 true 已归档 false 正常
	CreateTime int64                `bson:"create_time"` // 创建时间
	UpdateTime int64                `bson:"update_time"` // 修改时间
}
