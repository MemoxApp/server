package bce

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"time_speak_server/src/log"
	"time_speak_server/src/service/resource"
)

func Callback(conf Config, m *mongo.Collection) func(c *gin.Context) {
	if conf.CallbackToken == "" {
		panic("百度云回调token未配置")
	} else {
		log.Info("百度云回调地址: /notify/bce?token=" + conf.CallbackToken)
	}
	return func(c *gin.Context) {
		token := c.Query("token")
		if token != conf.CallbackToken {
			log.Error("百度云回调token错误")
			return
		}
		body := c.Request.Body
		bytes, err := io.ReadAll(body)
		if err != nil {
			log.Error("百度云回调请求读取失败", err)
			return
		}
		var callback CallbackRequest
		err = json.Unmarshal(bytes, &callback)
		if err != nil {
			log.Error("百度云回调请求解析失败", err)
			return
		}
		for _, event := range callback.Events {
			if event.EventType == "PutObject" || event.EventType == "PostObject" {
				if event.Content.Filesize > 0 && event.Content.Object != "" {
					size := event.Content.Filesize
					path := event.Content.Object
					_, err = resource.UpdateResourceSize(m, c.Request.Context(), path, int64(size))
					if err != nil {
						log.Error("百度云回调更新资源大小失败", err)
						return
					}
				}
			}
		}
		defer c.Request.Body.Close()
		c.String(200, "ok")
	}
}
