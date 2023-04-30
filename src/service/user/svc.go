package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"memox_server/src/exception"
	"memox_server/src/opts"
	"memox_server/src/service/cache"
	"strconv"
	"time"
)

type (
	Svc struct {
		Config
		m     *mongo.Collection
		redis *redis.Client
		c     *cache.Svc
	}
)

func NewUserSvc(conf Config, db *mongo.Database, r *redis.Client) *Svc {
	return &Svc{
		Config: conf,
		m:      db.Collection("user"),
		redis:  r,
		c:      cache.NewCacheSvc(r),
	}
}

// AddUser 添加用户 返回用户ID
func (s *Svc) AddUser(ctx context.Context, username, password, mail string) (primitive.ObjectID, error) {
	p, err := EncryptPassword(password)
	if err != nil {
		return primitive.NilObjectID, err
	}
	user := User{
		ObjectID:   primitive.NewObjectID(),
		Username:   username,
		Password:   p,
		Mail:       mail,
		CreateTime: time.Now().Unix(),
	}
	count, err := s.GetUserCount(ctx)
	if err != nil {
		return primitive.NilObjectID, err
	}
	if count == 0 {
		// 默认设置第一个用户为管理员
		user.Permission = 1
	}
	s.redis.Set(ctx, "UserCount", count+1, time.Minute*time.Duration(10))
	_, err = s.m.InsertOne(ctx, user)
	return user.ObjectID, err
}

// GetUser 获取用户
func (s *Svc) GetUser(ctx context.Context, id primitive.ObjectID) (u User, err error) {
	f := func() ([]byte, error) {
		user, err := s.getUser(ctx, id)
		if err != nil {
			return nil, err
		}
		return json.Marshal(user)
	}
	// Redis缓存
	result, err := s.c.Get(ctx, fmt.Sprintf("U-%s", id.Hex()), time.Minute*time.Duration(10), f)
	if err != nil {
		return
	}
	err = json.Unmarshal(result, &u)
	return
}

// GetUserCount 获取用户数量
func (s *Svc) GetUserCount(ctx context.Context) (c int64, err error) {
	f := func() ([]byte, error) {
		c, err = s.m.CountDocuments(ctx, bson.M{})
		if err != nil {
			return nil, err
		}
		return []byte(fmt.Sprintf("%d", c)), nil
	}
	// Redis缓存
	result, err := s.c.Get(ctx, "UserCount", time.Minute*time.Duration(10), f)
	if err != nil {
		return
	}
	count, err := strconv.Atoi(string(result))
	if err != nil {
		return
	}
	c = int64(count)
	return
}

// GetUserCountBySubscribe 获取指定订阅的用户数量
func (s *Svc) GetUserCountBySubscribe(ctx context.Context, id primitive.ObjectID) (c int64, err error) {
	c, err = s.m.CountDocuments(ctx, bson.M{"subscribe": id})
	return
}

// GetUserByMail 通过邮箱获取用户
func (s *Svc) GetUserByMail(ctx context.Context, mail string) (u User, err error) {
	err = s.m.FindOne(ctx, bson.M{"mail": mail}).Decode(&u)
	return
}

// UpdateUser 更新用户
func (s *Svc) UpdateUser(ctx context.Context, id primitive.ObjectID, opts ...opts.Option) error {
	toUpdate := bson.M{"profile_change_time": time.Now().Unix()}
	for _, f := range opts {
		toUpdate = f(toUpdate)
	}
	_, err := s.m.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": toUpdate})
	s.c.Del(ctx, fmt.Sprintf("U-%s", id.Hex()))
	return err
}

// getUser 通过id获取用户
func (s *Svc) getUser(ctx context.Context, id primitive.ObjectID) (u User, err error) {
	err = s.m.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	return
}

// CheckPasswordByMail 检查邮箱与密码是否匹配
func (s *Svc) CheckPasswordByMail(ctx context.Context, mail, password string) (bool, error) {
	u, err := s.GetUserByMail(ctx, mail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, exception.ErrMailNotFound
		}
		return false, err
	}
	return ComparePassword(u.Password, password)
}

// GetTokenByMail 生成Token
func (s *Svc) GetTokenByMail(ctx context.Context, mail string) (jwt JWTClaims, token string, err error) {
	var u User
	u, err = s.GetUserByMail(ctx, mail)
	if err != nil {
		return
	}
	jwt = JWTClaims{
		Subject:    "Memox",
		ExpiresAt:  time.Now().Add(time.Duration(s.TokenExpire) * time.Minute).Unix(),
		ID:         u.ID(),
		Permission: u.Permission,
	}
	token, err = GenerateJWTToken(jwt, s.TokenSecret)
	if err != nil {
		return
	}
	_, _ = s.m.UpdateOne(ctx, bson.M{"mail": mail}, bson.M{"$set": bson.M{"login_time": time.Now().Unix()}})
	return
}

// ParseToken 解析Token
func (s *Svc) ParseToken(token string) (*JWTClaims, error) {
	return ParseJWTToken(token, s.TokenSecret)
}
