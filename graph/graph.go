package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"time_speak_server/graph/generated"
	"time_speak_server/src/config"
	"time_speak_server/src/resolver"
)

func GraphqlHandler(conf config.Config, db *mongo.Database, redis *redis.Client) gin.HandlerFunc {

	c := generated.Config{
		Resolvers:  resolver.NewResolver(conf, db, redis),
		Directives: generated.DirectiveRoot{},
	}
	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// PlaygroundHandler Defining the Playground handler
func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
