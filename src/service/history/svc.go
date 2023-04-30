package history

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"time_speak_server/src/exception"
	"time_speak_server/src/service/memory"
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

func (s *Svc) NewHistory(ctx context.Context, oldMemory *memory.Memory) (string, error) {
	id, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return "", err
	}
	history := History{
		ObjectID:   primitive.NewObjectID(),
		Uid:        id,
		MemoryID:   oldMemory.ObjectID,
		Title:      oldMemory.Title,
		Content:    oldMemory.Content,
		HashTags:   oldMemory.HashTags,
		CreateTime: time.Now().Unix(),
	}
	_, err = s.m.InsertOne(ctx, history)
	return history.ObjectID.Hex(), err
}

// DeleteHistoryByMemoryID 根据memoryID删除历史
func (s *Svc) DeleteHistoryByMemoryID(ctx context.Context, memoryID primitive.ObjectID) error {
	_, err := s.m.DeleteMany(ctx, bson.M{"memory_id": memoryID})
	return err
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
	id, err := primitive.ObjectIDFromHex(memoryID)
	if err != nil {
		return nil, exception.ErrInvalidID
	}
	cursor, err := s.m.Find(ctx, bson.M{"memory_id": id}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &histories)
	return histories, err
}
