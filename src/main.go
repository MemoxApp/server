package main

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/gin-gonic/gin"
	"os"
	"time_speak_server/graph"
	"time_speak_server/src/config"
	"time_speak_server/src/db"
	"time_speak_server/src/log"
)

func main() {
	configFile, ok := os.LookupEnv("CONFIG_FILE")
	if !ok {
		log.Fatal("env CONFIG_FILE not set")
	}
	conf := config.MustReadConfigFile(configFile)

	if conf.Debug {
		log.SetDev()
	}

	database, redis, err := db.InitDB(conf)
	if err != nil {
		log.Fatal("fail to start server: initDB error")
		return
	}

	r := gin.Default()
	r.POST("/query", graph.GraphqlHandler(conf, database, redis))
	r.GET("/", graph.PlaygroundHandler())
	err = r.Run()
	if err != nil {
		panic(err)
	}
}
