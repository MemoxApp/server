package user

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"time_speak_server/src/exception"
)

// EncryptPassword 使用BCrypt加密密码
func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword 判断密码与加密值是否对应
func ComparePassword(hash string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// GenerateJWTToken 根据键值对生成token
func GenerateJWTToken(claims JWTClaims, secret string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// GetUserFromJwt 从token中获取用户信息
func GetUserFromJwt(c context.Context) (myID primitive.ObjectID, err error) {
	authInfo := GetJWTClaims(c)
	if authInfo.ID == "" {
		return primitive.NilObjectID, exception.ErrPermissionDenied
	}
	myID, err = primitive.ObjectIDFromHex(authInfo.ID)
	return
}

// GetJWTClaims 获得JWTClaims
func GetJWTClaims(c context.Context) Info {
	info, _ := graphql.GetOperationContext(c).Stats.GetExtension("Auth").(Info)
	return info
}

// ParseJWTToken 解析token
func ParseJWTToken(t, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(t, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return &JWTClaims{}, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return &JWTClaims{}, err
}
