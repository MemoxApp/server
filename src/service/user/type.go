package user

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户数据类型
type User struct {
	ObjectID           primitive.ObjectID `bson:"_id"`                  // 用户ID
	Username           string             `bson:"username"`             // 用户名
	Mail               string             `bson:"mail"`                 // 用户名
	Password           string             `bson:"password"`             // 密码（BCrypt）
	Avatar             string             `bson:"avatar"`               // 头像URL
	CreateTime         int64              `bson:"create_time"`          // 注册时间
	LoginTime          int64              `bson:"login_time"`           // 上次登录时间
	PasswordChangeTime int64              `bson:"password_change_time"` // 上次修改密码时间
	ProfileChangeTime  int64              `bson:"profile_change_time"`  // 上次修改资料时间
	Permission         int                `bson:"permission"`           // 权限 0 普通用户 1 管理员
	Subscribe          primitive.ObjectID `bson:"subscribe"`            // 订阅ID
}

func (User) IsSearchResult() {}

func (u User) ID() string {
	return u.ObjectID.Hex()
}

// JWTClaims JWT结构
type JWTClaims struct {
	Subject    string `json:"sub"`
	ID         string `json:"id"`
	Permission int    `json:"per"`
	ExpiresAt  int64  `json:"exp"`
}

// Valid 验证token是否有效
func (c JWTClaims) Valid() error {
	if jwt.TimeFunc().Unix() > c.ExpiresAt {
		return errors.New("token is expired")
	}
	return nil
}
