package main

//go:generate go get github.com/99designs/gqlgen@v0.17.30
//go:generate go run github.com/99designs/gqlgen generate

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"memox_server/graph"
	"memox_server/src/config"
	"memox_server/src/db"
	"memox_server/src/log"
	"memox_server/src/service/storage/bce"
	"os"
)

func main() {
	configFile, ok := os.LookupEnv("CONFIG_FILE")
	if !ok {
		log.Fatal("env CONFIG_FILE not set")
	}
	conf := config.MustReadConfigFile(configFile)
	if conf.Debug {
		log.SetDev()
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	database, redis, err := db.InitDB(conf)
	if err != nil {
		log.Fatal("启动服务失败，数据库连接失败")
		log.Fatal(err.Error())
		return
	}

	r := gin.Default()
	r.Any("/query", graph.GraphqlHandler(conf, database, redis))
	if conf.Debug {
		r.GET("/play", graph.PlaygroundHandler())
	}
	r.Static("/resources", conf.Storage.Local.Folder+"/resources")
	r.Any("/notify/bce", bce.Callback(conf.Storage.BCE, database.Collection("resource")))
	fmt.Printf(`
          _____                    _____                    _____                   _______                                 
         /\    \                  /\    \                  /\    \                 /::\    \                ______          
        /::\____\                /::\    \                /::\____\               /::::\    \              |::|   |         
       /::::|   |               /::::\    \              /::::|   |              /::::::\    \             |::|   |         
      /:::::|   |              /::::::\    \            /:::::|   |             /::::::::\    \            |::|   |         
     /::::::|   |             /:::/\:::\    \          /::::::|   |            /:::/~~\:::\    \           |::|   |         
    /:::/|::|   |            /:::/__\:::\    \        /:::/|::|   |           /:::/    \:::\    \          |::|   |         
   /:::/ |::|   |           /::::\   \:::\    \      /:::/ |::|   |          /:::/    / \:::\    \         |::|   |         
  /:::/  |::|___|______    /::::::\   \:::\    \    /:::/  |::|___|______   /:::/____/   \:::\____\        |::|   |         
 /:::/   |::::::::\    \  /:::/\:::\   \:::\    \  /:::/   |::::::::\    \ |:::|    |     |:::|    | ______|::|___|___ ____ 
/:::/    |:::::::::\____\/:::/__\:::\   \:::\____\/:::/    |:::::::::\____\|:::|____|     |:::|    ||:::::::::::::::::|    |
\::/    / ~~~~~/:::/    /\:::\   \:::\   \::/    /\::/    / ~~~~~/:::/    / \:::\    \   /:::/    / |:::::::::::::::::|____|
 \/____/      /:::/    /  \:::\   \:::\   \/____/  \/____/      /:::/    /   \:::\    \ /:::/    /   ~~~~~~|::|~~~|~~~      
             /:::/    /    \:::\   \:::\    \                  /:::/    /     \:::\    /:::/    /          |::|   |         
            /:::/    /      \:::\   \:::\____\                /:::/    /       \:::\__/:::/    /           |::|   |         
           /:::/    /        \:::\   \::/    /               /:::/    /         \::::::::/    /            |::|   |         
          /:::/    /          \:::\   \/____/               /:::/    /           \::::::/    /             |::|   |         
         /:::/    /            \:::\    \                  /:::/    /             \::::/    /              |::|   |         
        /:::/    /              \:::\____\                /:::/    /               \::/____/               |::|   |         
        \::/    /                \::/    /                \::/    /                 ~~                     |::|___|         
         \/____/                  \/____/                  \/____/                                          ~~

`)
	log.Info(fmt.Sprintf("Memox [Version] %s(%d) [Build] %d", config.VersionName, config.VersionCode, config.Build))
	err = r.Run()
	if err != nil {
		panic(err)
	}
}
