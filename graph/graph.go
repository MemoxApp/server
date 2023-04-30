package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"memox_server/graph/generated"
	"memox_server/src/config"
	"memox_server/src/resolver"
	"memox_server/src/service/user"
	"time"
)

func GraphqlHandler(conf config.Config, db *mongo.Database, redis *redis.Client) gin.HandlerFunc {
	c := generated.Config{
		Resolvers:  resolver.NewResolver(conf, db, redis),
		Directives: generated.DirectiveRoot{Auth: user.GqlAuth, Admin: user.GqlAdmin},
	}
	srv := handler.New(generated.NewExecutableSchema(c))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{
		MaxMemory:     32 * 1024 * 1024,
		MaxUploadSize: 20 * 1024 * 1024,
	})
	srv.SetQueryCache(lru.New(1000))
	srv.Use(&user.Extension{
		UserSvc: user.NewUserSvc(conf.User, db, redis),
	})

	if conf.Debug {
		srv.Use(extension.Introspection{})
	}

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(5))
		defer cancel()
		srv.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}

// PlaygroundHandler Defining the Playground handler
func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
