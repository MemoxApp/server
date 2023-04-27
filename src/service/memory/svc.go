package memory

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"time_speak_server/src/exception"
	"time_speak_server/src/service/cache"
	"time_speak_server/src/service/user"
)

type Svc struct {
	Config
	redis *redis.Client
	m     *mongo.Collection
	c     *cache.Svc
}
type Option func(bson.M) bson.M

func NewMemorySvc(conf Config, db *mongo.Database, redis *redis.Client) *Svc {
	return &Svc{
		Config: conf,
		redis:  redis,
		m:      db.Collection("memory"),
		c:      cache.NewCacheSvc(redis),
	}
}

func (s *Svc) CheckMemoryExist(ctx context.Context, title, content string) (bool, error) {
	var memory Memory
	err := s.m.FindOne(ctx, bson.M{"title": title, "content": content}).Decode(&memory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

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

func (s *Svc) UpdateMemory(ctx context.Context, id primitive.ObjectID, opts ...Option) error {
	toUpdate := bson.M{"update_time": time.Now().Unix()}
	for _, f := range opts {
		toUpdate = f(toUpdate)
	}
	_, err := s.m.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": toUpdate})
	s.c.Del(ctx, fmt.Sprintf("Memory-%s", id.Hex()))
	return err
}
func (s *Svc) DeleteMemory(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.m.DeleteOne(ctx, bson.M{"_id": id, "archived": true}) // 只有归档的才能删除
	s.c.Del(ctx, fmt.Sprintf("Memory-%s", id.Hex()))
	return err
}

func (s *Svc) GetMemory(ctx context.Context, id primitive.ObjectID) (*Memory, error) {
	f := func() ([]byte, error) {
		var memory Memory
		err := s.m.FindOne(ctx, bson.M{"_id": id}).Decode(&memory)
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

func WithTitle(t string) Option {
	return func(m bson.M) bson.M {
		m["title"] = t
		return m
	}
}

func WithArchived(t bool) Option {
	return func(m bson.M) bson.M {
		m["archived"] = t
		return m
	}
}

func WithContent(t string) Option {
	return func(m bson.M) bson.M {
		m["content"] = t
		return m
	}
}
func WithTags(t []primitive.ObjectID) Option {
	return func(m bson.M) bson.M {
		m["hash_tags"] = t
		return m
	}
}
