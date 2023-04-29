package main

//go:generate go get github.com/99designs/gqlgen@v0.17.30
//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/gin-gonic/gin"
	"os"
	"time_speak_server/graph"
	"time_speak_server/src/config"
	"time_speak_server/src/db"
	"time_speak_server/src/log"
	"time_speak_server/src/service/resource"
	"time_speak_server/src/service/storage"
	"time_speak_server/src/service/storage/bce"
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

	sto := storage.NewStorageSvc(conf.Storage, redis)
	resourceSvc := resource.NewResourceSvc(conf.Resource, database, redis, sto)

	r := gin.Default()
	r.Any("/query", graph.GraphqlHandler(conf, database, redis))
	r.GET("/", graph.PlaygroundHandler())
	r.Static("/resources", conf.Storage.Local.Folder+"/resources")
	r.POST("/notify/bce", bce.Callback(resourceSvc))
	err = r.Run()
	if err != nil {
		panic(err)
	}
}
