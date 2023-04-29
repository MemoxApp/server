package memory

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

// NewMemorySvc 创建记忆服务
func NewMemorySvc(conf Config, db *mongo.Database, redis *redis.Client) *Svc {
	return &Svc{
		Config: conf,
		redis:  redis,
		m:      db.Collection("memory"),
		c:      cache.NewCacheSvc(redis),
	}
}

// CheckMemoryExist 检查记忆是否存在
func (s *Svc) CheckMemoryExist(ctx context.Context, title, content string) (bool, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return false, err
	}
	var memory Memory
	err = s.m.FindOne(ctx, bson.M{"uid": uid, "title": title, "content": content}).Decode(&memory) // 检查自己的记忆
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// NewMemory 创建记忆
func (s *Svc) NewMemory(ctx context.Context, title, content string, tags []primitive.ObjectID) (string, error) {
	exist, err := s.CheckMemoryExist(ctx, title, content)
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
	memory := Memory{
		ObjectID:   primitive.NewObjectID(),
		Uid:        id,
		Title:      title,
		Content:    content,
		HashTags:   tags,
		Archived:   false,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	_, err = s.m.InsertOne(ctx, memory)
	return memory.ObjectID.Hex(), err
}

// UpdateMemory 更新记忆
func (s *Svc) UpdateMemory(ctx context.Context, id primitive.ObjectID, opts ...opts.Option) error {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return err
	}
	toUpdate := bson.M{"update_time": time.Now().Unix()}
	for _, f := range opts {
		toUpdate = f(toUpdate)
	}
	_, err = s.m.UpdateOne(ctx, bson.M{"uid": uid, "_id": id}, bson.M{"$set": toUpdate}) // 只能更新自己的记忆
	s.c.Del(ctx, fmt.Sprintf("Memory-%s", id.Hex()))
	return err
}

// DeleteMemory 删除记忆
func (s *Svc) DeleteMemory(ctx context.Context, id primitive.ObjectID) error {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return err
	}
	_, err = s.m.DeleteOne(ctx, bson.M{"uid": uid, "_id": id, "archived": true}) // 只能删除自己的已归档的记忆
	s.c.Del(ctx, fmt.Sprintf("Memory-%s", id.Hex()))
	return err
}

// GetMemory 获取记忆
func (s *Svc) GetMemory(ctx context.Context, id primitive.ObjectID) (*Memory, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	f := func() ([]byte, error) {
		var memory Memory
		err := s.m.FindOne(ctx, bson.M{"uid": uid, "_id": id}).Decode(&memory) // 只能获取自己的记忆
		if err != nil {
			return nil, err
		}
		return bson.Marshal(memory)
	}
	var memory Memory
	bytes, err := s.c.Get(ctx, fmt.Sprintf("Memory-%s", id.Hex()), time.Hour, f)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(bytes, &memory)
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

// GetMemories 获取记忆列表
func (s *Svc) GetMemories(ctx context.Context, page, size int64, byCreate, desc, archived bool) ([]*Memory, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	var memory []*Memory
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
	data, err := s.m.Find(ctx, bson.M{"uid": uid, "archived": archived}, &options.FindOptions{
		Skip:  &skip,
		Limit: &size,
		Sort:  sort,
	})
	if err != nil {
		return nil, err
	}
	defer data.Close(ctx)
	err = data.All(ctx, &memory)
	if err != nil {
		return nil, err
	}
	return memory, nil
}

// GetMemoriesByHashTag 根据标签获取记忆列表
func (s *Svc) GetMemoriesByHashTag(ctx context.Context, tagID primitive.ObjectID, page, size int64, byCreate, desc, archived bool) ([]*Memory, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	var memory []*Memory
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
	data, err := s.m.Find(ctx, bson.M{"uid": uid, "archived": archived, "hash_tags": bson.M{"$in": []primitive.ObjectID{tagID}}}, &options.FindOptions{
		Skip:  &skip,
		Limit: &size,
		Sort:  sort,
	})
	if err != nil {
		return nil, err
	}
	defer data.Close(ctx)
	err = data.All(ctx, &memory)
	if err != nil {
		return nil, err
	}
	return memory, nil
}

// GetMemoriesCountByHashTag 根据标签获取记忆数量
func (s *Svc) GetMemoriesCountByHashTag(ctx context.Context, tagID primitive.ObjectID) (int64, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return 0, err
	}
	data, err := s.m.CountDocuments(ctx,
		bson.M{
			"uid":       uid,
			"hash_tags": bson.M{"$in": []primitive.ObjectID{tagID}}},
	)
	if err != nil {
		return 0, err
	}
	return data, nil
}
