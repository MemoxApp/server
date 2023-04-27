package history

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s *Svc) GetHistories(ctx context.Context, memoryID string, page, size int64, desc bool) ([]*History, error) {
	order := 1
	if desc {
		order = -1
	}
	skip := page * size
	var histories []*History
	opts := &options.FindOptions{
		Skip:  &skip,
		Limit: &size,
		Sort:  bson.M{"create_time": order},
	}
	cursor, err := s.m.Find(ctx, bson.M{"memory_id": memoryID}, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &histories)
	return histories, err
}
