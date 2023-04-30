package comment

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"time_speak_server/src/exception"
	"time_speak_server/src/opts"
	"time_speak_server/src/service/cache"
	"time_speak_server/src/service/user"
)

type Svc struct {
	Config
	redis *redis.Client
	m     *mongo.Collection
	c     *cache.Svc
}

// NewCommentSvc 创建回复服务
func NewCommentSvc(conf Config, db *mongo.Database, redis *redis.Client) *Svc {
	return &Svc{
		Config: conf,
		redis:  redis,
		m:      db.Collection("comment"),
		c:      cache.NewCacheSvc(redis),
	}
}

// CheckCommentExist 检查回复是否存在
func (s *Svc) CheckCommentExist(ctx context.Context, content string, commentID, parentID primitive.ObjectID) (bool, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return false, err
	}
	var comment Comment
	err = s.m.FindOne(ctx, bson.M{"uid": uid, "content": content, "comment_id": commentID, "parent_id": parentID}).Decode(&comment) // 检查自己的回复
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// NewComment 创建回复
func (s *Svc) NewComment(ctx context.Context, content string, commentID, parentID primitive.ObjectID, tags []primitive.ObjectID) (string, error) {
	exist, err := s.CheckCommentExist(ctx, content, commentID, parentID)
	if err != nil {
		return "", err
	}
	if exist {
		return "", exception.ErrContentExist
	}
	id, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return "", err
	}
	comment := Comment{
		ObjectID:   primitive.NewObjectID(),
		Uid:        id,
		Content:    content,
		CommentID:  commentID,
		ParentID:   parentID,
		HashTags:   tags,
		Archived:   false,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	_, err = s.m.InsertOne(ctx, comment)
	return comment.ObjectID.Hex(), err
}

// UpdateComment 更新回复
func (s *Svc) UpdateComment(ctx context.Context, id primitive.ObjectID, opts ...opts.Option) error {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return err
	}
	toUpdate := bson.M{"update_time": time.Now().Unix()}
	for _, f := range opts {
		toUpdate = f(toUpdate)
	}
	_, err = s.m.UpdateOne(ctx, bson.M{"uid": uid, "_id": id}, bson.M{"$set": toUpdate}) // 只能更新自己的回复
	s.c.Del(ctx, fmt.Sprintf("Comment-%s", id.Hex()))
	return err
}

// DeleteComment 删除回复
func (s *Svc) DeleteComment(ctx context.Context, id primitive.ObjectID) error {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return err
	}
	comment, err := s.GetComment(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return exception.ErrCommentNotFound
		}
		return err
	}
	if !comment.Archived {
		return exception.ErrCommentNotArchived
	}
	result, err := s.m.DeleteOne(ctx, bson.M{"_id": id, "uid": uid, "archived": true}) // 只有归档的才能删除
	s.c.Del(ctx, fmt.Sprintf("Comment-%s", id.Hex()))
	if result.DeletedCount == 0 {
		return exception.ErrCommentNotFound
	}
	return err
}

// GetComment 获取回复
func (s *Svc) GetComment(ctx context.Context, id primitive.ObjectID) (*Comment, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	var comment Comment
	err = s.m.FindOne(ctx, bson.M{"uid": uid, "_id": id}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.ErrCommentNotFound
		}
		return nil, err
	}
	return &comment, nil
}

// GetComments 获取回复列表
func (s *Svc) GetComments(ctx context.Context, parentID string, page, size int64, byCreate, desc, archived bool) ([]*Comment, error) {
	var comment []*Comment
	skip := page * size
	order := 1
	if desc {
		order = -1
	}
	sort := bson.M{
		"update_time": order,
	}
	if byCreate {
		sort = bson.M{
			"create_time": order,
		}
	}
	id, err := primitive.ObjectIDFromHex(parentID)
	if err != nil {
		return nil, err
	}
	data, err := s.m.Find(ctx, bson.M{"parent_id": id, "archived": archived}, &options.FindOptions{
		Skip:  &skip,
		Limit: &size,
		Sort:  sort,
	})
	if err != nil {
		return nil, err
	}
	defer data.Close(ctx)
	err = data.All(ctx, &comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// GetArchivedComments 获取已归档回复列表
func (s *Svc) GetArchivedComments(ctx context.Context, page, size int64, byCreate, desc bool) ([]*Comment, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	var comment []*Comment
	skip := page * size
	order := 1
	if desc {
		order = -1
	}
	sort := bson.M{
		"update_time": order,
	}
	if byCreate {
		sort = bson.M{
			"create_time": order,
		}
	}
	data, err := s.m.Find(ctx, bson.M{"uid": uid, "archived": true}, &options.FindOptions{
		Skip:  &skip,
		Limit: &size,
		Sort:  sort,
	}) // 只能获取自己的归档回复
	if err != nil {
		return nil, err
	}
	defer data.Close(ctx)
	err = data.All(ctx, &comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}
