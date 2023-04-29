package bce

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"time_speak_server/src/log"
	"time_speak_server/src/service/resource"
)

func Callback(svc *resource.Svc) func(c *gin.Context) {
	return func(c *gin.Context) {
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
		if callback.Events.EventType == "PutObject" || callback.Events.EventType == "PostObject" {
			if callback.Events.Content.Filesize > 0 && callback.Events.Content.Object != "" {
				size := callback.Events.Content.Filesize
				path := callback.Events.Content.Object
				err = svc.UpdateResourceSize(c.Request.Context(), path, int64(size))
				if err != nil {
					log.Error("百度云回调更新资源大小失败", err)
					return
				}
			}
		}
		defer c.Request.Body.Close()
		c.String(200, "ok")
	}
}
