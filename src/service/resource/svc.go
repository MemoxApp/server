package resource

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"memox_server/src/exception"
	"memox_server/src/opts"
	"memox_server/src/service/cache"
	"memox_server/src/service/storage/local"
	"memox_server/src/service/storage/utils"
	"memox_server/src/service/user"
	"strings"
	"time"
)

type Svc struct {
	Config
	redis *redis.Client
	m     *mongo.Collection
	c     *cache.Svc
	sto   utils.Service
}

// NewResourceSvc 创建资源服务
func NewResourceSvc(conf Config, db *mongo.Database, redis *redis.Client, sto utils.Service) *Svc {
	return &Svc{
		Config: conf,
		redis:  redis,
		m:      db.Collection("resource"),
		c:      cache.NewCacheSvc(redis),
		sto:    sto,
	}
}

// CheckResourceExist 检查资源是否存在
func (s *Svc) CheckResourceExist(ctx context.Context, path string) (*Resource, error) {
	id, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return nil, err
	}
	var resource Resource
	err = s.m.FindOne(ctx, bson.M{"path": path, "uid": id}).Decode(&resource)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &resource, nil
}

// NewResource 创建资源
func (s *Svc) NewResource(ctx context.Context, path string, size int64) (string, error) {
	id, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return "", err
	}
	exist, err := s.CheckResourceExist(ctx, path)
	if err != nil {
		return "", err
	}
	if exist != nil {
		if exist.Size > 0 {
			return exist.ObjectID.Hex(), exception.ErrResourceExist
		} else {
			// 创建了但未使用的资源
			return exist.ObjectID.Hex(), nil
		}
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
	toUpdate := bson.M{"update_time": time.Now().Unix()}
	for _, f := range opts {
		toUpdate = f(toUpdate)
	}
	_, err := s.m.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": toUpdate})
	s.c.Del(ctx, fmt.Sprintf("Resource-%s", id.Hex()))
	return err
}

// UpdateResourceReferences 更新资源引用
func (s *Svc) UpdateResourceReferences(ctx context.Context, ref, memoryID string) error {
	var resID primitive.ObjectID
	var refArray []primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(ref)
	if err != nil {
		return exception.ErrInvalidID
	}
	res, err := s.GetResource(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return exception.ErrResourceNotFound
		}
		return err
	} else {
		resID = res.ObjectID
		refArray = res.Ref
	}
	memory, err := primitive.ObjectIDFromHex(memoryID)
	if err != nil {
		return err
	}
	refArray = append(refArray, memory)
	refArray = UniqueArr(refArray) // 引用去重
	err = s.UpdateResource(ctx, resID, opts.With("ref", refArray))
	if err != nil {
		return err
	}
	return nil
}

// UpdateResourceSize 插入或更新资源
func UpdateResourceSize(m *mongo.Collection, ctx context.Context, fileName string, size int64) (string, error) {
	res, err := GetResourceByPath(m, ctx, fileName)
	if err != nil {
		return "", err
	}
	err = UpdateResource(m, ctx, res.ObjectID, opts.With("size", size))
	if err != nil {
		return "", err
	}
	return res.ObjectID.Hex(), nil
}

// GetResourceByPath 通过路径获取资源
func GetResourceByPath(m *mongo.Collection, ctx context.Context, path string) (*Resource, error) {
	split := strings.Split(path, "/")
	if len(split) < 3 {
		return nil, exception.ErrInvalidPath
	}
	userId := split[1]
	uid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, exception.ErrInvalidID
	}
	fileName := split[2]
	var resource Resource
	err = m.FindOne(ctx, bson.M{"uid": uid, "path": fileName}).Decode(&resource) // 只能获取自己的资源
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// UpdateResource 更新资源
func UpdateResource(m *mongo.Collection, ctx context.Context, id primitive.ObjectID, opts ...opts.Option) error {
	toUpdate := bson.M{"update_time": time.Now().Unix()}
	for _, f := range opts {
		toUpdate = f(toUpdate)
	}
	_, err := m.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": toUpdate})
	return err
}

func (s *Svc) UpdateReferences(ctx context.Context, oldContent, newContent string, memoryID string) error {
	oldReferences := utils.PickupReferences(oldContent)
	newReferences := utils.PickupReferences(newContent)
	fmt.Printf("old: %v, new: %v\n", oldReferences, newReferences)
	removeReferences, insertReferences := DiffArray(oldReferences, newReferences)
	fmt.Printf("remove: %v, insert: %v\n", removeReferences, insertReferences)
	if len(insertReferences) > 0 {
		for _, ref := range insertReferences {
			err := s.UpdateResourceReferences(ctx, ref, memoryID)
			if err != nil {
				return err
			}
		}
	}
	if len(removeReferences) > 0 {
		for _, ref := range removeReferences {
			err := s.RemoveResourceReference(ctx, ref, memoryID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RemoveResourceReference 移除资源引用
func (s *Svc) RemoveResourceReference(ctx context.Context, ref, memoryID string) error {
	var refArray []primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(ref)
	if err != nil {
		return exception.ErrInvalidID
	}
	res, err := s.GetResource(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	refArray = res.Ref
	memory, err := primitive.ObjectIDFromHex(memoryID)
	if err != nil {
		return err
	}
	refArray = RemoveFromArr(refArray, memory) // 移除引用
	err = s.UpdateResource(ctx, res.ObjectID, opts.With("ref", refArray))
	if err != nil {
		return err
	}
	return nil
}

// DeleteResource 删除资源
func (s *Svc) DeleteResource(ctx context.Context, id primitive.ObjectID) error {
	uid, err := user.GetUserFromJwt(ctx)
	if err != nil {
		return err
	}
	_, err = s.m.DeleteOne(ctx, bson.M{"uid": uid, "_id": id}) // 只能删除自己的资源
	if err != nil {
		return err
	}
	// 删除资源
	suc, err := s.sto.Delete(ctx, id.Hex())
	if err != nil {
		return err
	}
	if !suc {
		return exception.ErrDeleteResource
	}
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
		if err == mongo.ErrNoDocuments {
			return nil, exception.ErrResourceNotFound
		}
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

func (s *Svc) InsertResourceUrl(ctx context.Context, content string) (string, error) {
	refs := utils.PickupReferences(content)
	for _, ref := range refs {
		id, err := primitive.ObjectIDFromHex(ref)
		if err != nil {
			return "", exception.ErrInvalidID
		}
		res, err := s.GetResource(ctx, id)
		if err != nil {
			return "", err
		}
		url, err := s.sto.GetUrl(ctx, res.Path)
		if err != nil {
			return "", err
		}
		content = strings.ReplaceAll(content, "${"+ref+"}", url)
	}
	return content, nil
}

func (s *Svc) LocalUpload(ctx context.Context, session string, upload graphql.Upload) (string, error) {
	localUploader, ok := s.sto.(*local.Local)
	if !ok {
		return "", exception.ErrInvalidStorageProvider
	}
	l, size, err := localUploader.Upload(ctx, session, upload)
	if err != nil {
		return "", err
	}
	println("fileName:" + l)
	id, err := UpdateResourceSize(s.m, ctx, l, size)
	return id, err
}
