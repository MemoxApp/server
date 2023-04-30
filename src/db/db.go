package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"memox_server/src/config"
	"memox_server/src/log"
	"time"
)

func InitDB(conf config.Config) (db *mongo.Database, r *redis.Client, err error) {
	db, err = initMongo(conf.DB.MongoAddr, conf.DB.MongoDB)
	if err != nil {
		return
	}
	r, err = initRedis(conf.DB.RedisAddr, conf.DB.RedisDB)
	return
}

// initMongo 连接Mongo数据库
func initMongo(address string, database string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s", address)
	mongoC, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Error("fail to connect mongo",
			"addr", address, "err", err)
		return nil, err
	}

	err = mongoC.Ping(ctx, nil)
	if err != nil {
		log.Error("fail to ping mongo",
			"addr", address, "err", err)
		return nil, err
	}

	log.Info("successfully init mongo",
		"addr", address)
	return mongoC.Database(database), nil
}

// initRedis 连接Redis数据库
func initRedis(address string, database int) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	r := redis.NewClient(
		&redis.Options{
			Addr: address,
			DB:   database,
		},
	)

	err := r.Ping(ctx).Err()
	if err != nil {
		log.Error("fail to ping redis",
			"addr", address, "err", err)
		return nil, err
	}

	log.Info("successfully init redis",
		"addr", address)
	return r, nil
}
