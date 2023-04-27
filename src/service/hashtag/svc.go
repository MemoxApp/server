package hashtag

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
	"time"
	"time_speak_server/src/service/cache"
	"time_speak_server/src/service/user"
)

type Svc struct {
	Config
	redis *redis.Client
	c     *cache.Svc
	m     *mongo.Collection
}

func NewHashtagSvc(conf Config, db *mongo.Database, redis *redis.Client) *Svc {
	return &Svc{
		Config: conf,
		redis:  redis,
		m:      db.Collection("hashtag"),
		c:      cache.NewCacheSvc(redis),
	}
}

func (s *Svc) NewHashtag(ctx context.Context, name string) (primitive.ObjectID, error) {
	id, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return primitive.NilObjectID, err
	}
	hashtag := HashTag{
		ObjectID:   primitive.NewObjectID(),
		Uid:        id,
		Name:       name,
		Archived:   false,
		CreateTime: time.Now().Unix(),
	}
	_, err = s.m.InsertOne(ctx, hashtag)
	return hashtag.ObjectID, err
}

func (s *Svc) GetOrInsertHashtag(ctx context.Context, name string) (primitive.ObjectID, error) {
	var tag HashTag
	err := s.m.FindOne(ctx, bson.M{"name": name}).Decode(&tag)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			hashtag, err := s.NewHashtag(ctx, name)
			if err != nil {
				return primitive.NilObjectID, err
			}
			return hashtag, nil
		}
		return primitive.NilObjectID, err
	}
	return tag.ObjectID, nil
}

func (s *Svc) GetHashTag(ctx context.Context, name string) (*HashTag, error) {
	f := func() ([]byte, error) {
		var tag HashTag
		err := s.m.FindOne(ctx, bson.M{"name": name}).Decode(&tag)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errHashTagNotFound
			}
			return nil, err
		}
		return bson.Marshal(tag)
	}
	var tag HashTag
	// Redis 缓存
	result, err := s.c.Get(ctx, fmt.Sprintf("Tag-%s", name), time.Minute*time.Duration(10), f)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(result, &tag)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *Svc) GetHashTagByID(ctx context.Context, id primitive.ObjectID) (*HashTag, error) {
	f := func() ([]byte, error) {
		var tag HashTag
		err := s.m.FindOne(ctx, bson.M{"_id": id}).Decode(&tag)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, errHashTagNotFound
			}
			return nil, err
		}
		return bson.Marshal(tag)
	}
	var tag HashTag
	// Redis 缓存
	result, err := s.c.Get(ctx, fmt.Sprintf("#-%s", id.Hex()), time.Minute*time.Duration(10), f)
	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(result, &tag)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *Svc) MakeHashTags(ctx context.Context, content string) ([]primitive.ObjectID, error) {
	tags := ParseHashTags(content)
	var ids []primitive.ObjectID
	for _, tag := range tags {
		f := func() ([]byte, error) {
			hashtag, err := s.GetOrInsertHashtag(ctx, tag)
			if err != nil {
				return nil, err
			}
			return []byte(hashtag.Hex()), nil
		}
		// Redis 缓存
		result, err := s.c.Get(ctx, fmt.Sprintf("Tag-%s", tag), time.Minute*time.Duration(10), f)
		if err != nil {
			return nil, err
		}
		id, err := primitive.ObjectIDFromHex(string(result))
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func ParseHashTags(content string) []string {
	r, _ := regexp.Compile("#\\S[^\\n]+? ")
	return r.FindAllString(content, -1)
}
