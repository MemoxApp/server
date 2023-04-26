package memory

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"time_speak_server/src/exception"
	"time_speak_server/src/service/user"
)

type Svc struct {
	Config
	redis *redis.Client
	m     *mongo.Collection
}

func NewMemorySvc(conf Config, db *mongo.Database, redis *redis.Client) *Svc {
	return &Svc{
		Config: conf,
		redis:  redis,
		m:      db.Collection("memory"),
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
