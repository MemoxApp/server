package history

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"time_speak_server/src/service/user"
)

type Svc struct {
	Config
	redis *redis.Client
	m     *mongo.Collection
}

func NewHistorySvc(conf Config, db *mongo.Database, redis *redis.Client) *Svc {
	return &Svc{
		Config: conf,
		redis:  redis,
		m:      db.Collection("history"),
	}
}

func (s *Svc) NewHistory(ctx context.Context, title, content string, tags []primitive.ObjectID) (string, error) {
	id, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return "", err
	}
	history := History{
		ObjectID:   primitive.NewObjectID(),
		Uid:        id,
		Title:      title,
		Content:    content,
		HashTags:   tags,
		CreateTime: time.Now().Unix(),
	}
	_, err = s.m.InsertOne(ctx, history)
	return history.ObjectID.Hex(), err
}
