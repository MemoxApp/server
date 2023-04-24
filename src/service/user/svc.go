package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"time_speak_server/src/service/cache"
)

type (
	Svc struct {
		Config
		m     *mongo.Collection
		redis *redis.Client
		c     cache.Svc
	}
	Option func(bson.M) bson.M
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
func (s *Svc) AddUser(ctx context.Context, username, password string) (primitive.ObjectID, error) {
	p, err := EncryptPassword(password)
	if err != nil {
		return primitive.NilObjectID, err
	}

	user := User{
		ObjectID:   primitive.NewObjectID(),
		Username:   username,
		Password:   p,
		CreateTime: time.Now().Unix(),
	}

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

// GetUserByUsername 通过用户名获取用户
func (s *Svc) GetUserByUsername(ctx context.Context, username string) (u User, err error) {
	err = s.m.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	return
}

// UpdateUser 更新用户
func (s *Svc) UpdateUser(ctx context.Context, id primitive.ObjectID, opts ...Option) error {
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

// CheckPasswordByUsername 检查用户名与密码是否匹配
func (s *Svc) CheckPasswordByUsername(ctx context.Context, username, password string) (bool, error) {
	u, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return ComparePassword(u.Password, password)
}

// GetTokenByUsername 生成Token
func (s *Svc) GetTokenByUsername(ctx context.Context, username string) (jwt JWTClaims, token string, err error) {
	var u User
	u, err = s.GetUserByUsername(ctx, username)
	if err != nil {
		return
	}
	jwt = JWTClaims{
		Subject:    "TimeSpeak",
		ExpiresAt:  time.Now().Add(time.Duration(s.TokenExpire) * time.Minute).Unix(),
		ID:         u.ID(),
		Permission: u.Permission,
	}
	token, err = GenerateJWTToken(jwt, s.TokenSecret)
	if err != nil {
		return
	}
	_, _ = s.m.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": bson.M{"login_time": time.Now().Unix()}})
	return
}

// ParseToken 解析Token
func (s *Svc) ParseToken(ctx context.Context, token string) (*JWTClaims, error) {
	return ParseJWTToken(token, s.TokenSecret)
}

// UpdateUserPassword 更新用户密码
func (s *Svc) UpdateUserPassword(ctx context.Context, username string, password string) error {
	p, err := EncryptPassword(password)
	if err != nil {
		return err
	}
	_, err = s.m.UpdateOne(ctx, bson.M{"username": username},
		bson.M{"$set": bson.M{"password": p, "password_change_time": time.Now().Unix()}})
	// 这里不需要处理缓存，因为登录也是用username获取用户的，缓存只做了通过id获取
	return err
}
