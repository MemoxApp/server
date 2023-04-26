package user

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"strings"
	"time_speak_server/src/exception"
	"time_speak_server/src/log"
)

type Info struct {
	ID         string
	Permission int
}

// Extension 验证token并加入到context中
type Extension struct {
	UserSvc *Svc
}

var _ interface {
	graphql.OperationParameterMutator
	graphql.HandlerExtension
} = &Extension{}

func (a *Extension) ExtensionName() string {
	return "AuthExtension"
}

func (a *Extension) Validate(schema graphql.ExecutableSchema) error {
	return nil
}

func (a *Extension) MutateOperationParameters(ctx context.Context, rawParams *graphql.RawParams) *gqlerror.Error {
	jwtStr := rawParams.Headers.Get("Authorization")

	var err error
	claims := &JWTClaims{}
	if jwtStr != "" {
		jwtStr = strings.TrimPrefix(jwtStr, "Bearer ")
		claims, err = a.UserSvc.ParseToken(jwtStr)
		if err != nil {
			log.Warn("fail to parse jwt token", "token", jwtStr, "err", err)
		}
	}
	//ip := rawParams.Headers.Get("X-Real-IP")
	graphql.GetOperationContext(ctx).Stats.SetExtension("Auth", Info{
		ID:         claims.ID,
		Permission: claims.Permission,
	})
	return nil
}

// GqlAuth 鉴权
func GqlAuth(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	_, ok := graphql.GetOperationContext(ctx).Stats.GetExtension("Auth").(Info)
	if ok {
		return next(ctx)
	} else {
		return nil, exception.ErrPermissionDenied
	}
}
