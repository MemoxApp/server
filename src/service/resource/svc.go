package resource

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

// NewResourceSvc 创建资源服务
func NewResourceSvc(conf Config, db *mongo.Database, redis *redis.Client) *Svc {
	return &Svc{
		Config: conf,
		redis:  redis,
		m:      db.Collection("resource"),
		c:      cache.NewCacheSvc(redis),
	}
}

// CheckResourceExist 检查资源是否存在
func (s *Svc) CheckResourceExist(ctx context.Context, path string) (bool, error) {
	var resource Resource
	err := s.m.FindOne(ctx, bson.M{"path": path}).Decode(&resource)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return true, err
	}
	return true, nil
}

// NewResource 创建资源
func (s *Svc) NewResource(ctx context.Context, path string, size int64) (string, error) {
	exist, err := s.CheckResourceExist(ctx, path)
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
	resource := Resource{
		ObjectID:   primitive.NewObjectID(),
		Uid:        id,
		Path:       path,
		Size:       size,
		Ref:        nil,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}
	_, err = s.m.InsertOne(ctx, resource)
	return resource.ObjectID.Hex(), err
}

// UpdateResource 更新资源
func (s *Svc) UpdateResource(ctx context.Context, id primitive.ObjectID, opts ...opts.Option) error {
	// todo check token
	toUpdate := bson.M{"update_time": time.Now().Unix()}
	for _, f := range opts {
		toUpdate = f(toUpdate)
	}
	_, err := s.m.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": toUpdate})
	s.c.Del(ctx, fmt.Sprintf("Resource-%s", id.Hex()))
	return err
}

// DeleteResource 删除资源
func (s *Svc) DeleteResource(ctx context.Context, id primitive.ObjectID) error {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return err
	}
	_, err = s.m.DeleteOne(ctx, bson.M{"uid": uid, "_id": id, "ref": nil}) // 只有没有引用的资源才能删除 // 只能删除自己的资源
	// todo DELETE 删除文件逻辑
	s.c.Del(ctx, fmt.Sprintf("Resource-%s", id.Hex()))
	return err
}

// GetResource 获取资源
func (s *Svc) GetResource(ctx context.Context, id primitive.ObjectID) (*Resource, error) {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	var resource Resource
	err = s.m.FindOne(ctx, bson.M{"uid": uid, "_id": id}).Decode(&resource) // 只能获取自己的资源
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetResourceUsed 获取资源使用情况
func (s *Svc) GetResourceUsed(ctx context.Context, uid primitive.ObjectID) (int64, error) {
	var resource []*Resource
	cursor, err := s.m.Find(ctx, bson.M{"uid": uid}) // 获取用户的所有资源
	if err != nil {
		return -1, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &resource)
	if err != nil {
		return -1, err
	}
	var size int64
	for _, v := range resource {
		size += v.Size
	}
	return size, nil
}

// GetResources 获取资源列表
func (s *Svc) GetResources(ctx context.Context, uid primitive.ObjectID, page, size int64, byCreate, desc bool) ([]*Resource, error) {
	var resource []*Resource
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
	data, err := s.m.Find(ctx, bson.M{"uid": uid}, &options.FindOptions{
		Skip:  &skip,
		Limit: &size,
		Sort:  sort,
	}) // 只能获取自己的资源
	if err != nil {
		return nil, err
	}
	defer data.Close(ctx)
	err = data.All(ctx, &resource)
	if err != nil {
		return nil, err
	}
	return resource, nil
}
