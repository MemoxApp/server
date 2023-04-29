package subscribe

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
)

type Svc struct {
	Config
	redis *redis.Client
	m     *mongo.Collection
	c     *cache.Svc
}

// NewSubscribeSvc 创建订阅服务
func NewSubscribeSvc(conf Config, db *mongo.Database, redis *redis.Client) *Svc {
	return &Svc{
		Config: conf,
		redis:  redis,
		m:      db.Collection("subscribe"),
		c:      cache.NewCacheSvc(redis),
	}
}

// CheckSubscribeExist 检查订阅是否存在
func (s *Svc) CheckSubscribeExist(ctx context.Context, name string) (bool, error) {
	var subscribe Subscribe
	err := s.m.FindOne(ctx, bson.M{"name": name}).Decode(&subscribe)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// NewSubscribe 创建订阅
func (s *Svc) NewSubscribe(ctx context.Context, name string, capacity int64, enable bool) (string, error) {
	exist, err := s.CheckSubscribeExist(ctx, name)
	if err != nil {
		return "", err
	}
	if exist {
		return "", exception.ErrContentExist
	}
	subscribe := Subscribe{
		ObjectID:   primitive.NewObjectID(),
		Name:       name,
		Capacity:   capacity,
		Enabled:    enable,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	_, err = s.m.InsertOne(ctx, subscribe)
	return subscribe.ObjectID.Hex(), err
}

// UpdateSubscribe 更新订阅
func (s *Svc) UpdateSubscribe(ctx context.Context, id primitive.ObjectID, opts ...opts.Option) error {
	toUpdate := bson.M{"update_time": time.Now().Unix()}
	for _, f := range opts {
		toUpdate = f(toUpdate)
	}
	_, err := s.m.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": toUpdate})
	s.c.Del(ctx, fmt.Sprintf("Subscribe-%s", id.Hex()))
	return err
}

// DeleteSubscribe 删除订阅
func (s *Svc) DeleteSubscribe(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.m.DeleteOne(ctx, bson.M{"_id": id, "enabled": true}) // 只有正常状态的订阅才能删除
	s.c.Del(ctx, fmt.Sprintf("Subscribe-%s", id.Hex()))
	return err
}

// GetSubscribe 获取订阅
func (s *Svc) GetSubscribe(ctx context.Context, id primitive.ObjectID) (*Subscribe, error) {
	var subscribe Subscribe
	if id == primitive.NilObjectID {
		return &Subscribe{
			ObjectID:   primitive.NilObjectID,
			Name:       s.DefaultSubscribeName,
			Capacity:   s.DefaultCapacity,
			Enabled:    true,
			CreateTime: 0,
			UpdateTime: 0,
		}, nil
	}
	err := s.m.FindOne(ctx, bson.M{"_id": id}).Decode(&subscribe)
	if err != nil {
		return nil, err
	}
	return &subscribe, nil
}

// GetSubscribes 获取订阅列表
func (s *Svc) GetSubscribes(ctx context.Context, enabled bool) ([]*Subscribe, error) {
	var subscribe []*Subscribe
	sort := bson.M{
		"capacity": 1, // 按容量升序
	}
	data, err := s.m.Find(ctx, bson.M{"enabled": enabled}, &options.FindOptions{
		Sort: sort,
	})
	if err != nil {
		return nil, err
	}
	defer data.Close(ctx)
	err = data.All(ctx, &subscribe)
	if err != nil {
		return nil, err
	}
	return subscribe, nil
}
