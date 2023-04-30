package bce

import (
	"fmt"
	"time"
	"time_speak_server/src/service/storage/utils"
)

func (b *BCE) getWritePermissionACL(userID, filename string) string {
	return ` {
  "accessControlList": [
   {   
    "service":"bce:bos",
    "region":"` + b.Config.Region + `",
    "effect": "Allow",
    "resource": ["` + utils.GenerateResourcePath(userID, filename) + `"],
    "permission": ["WRITE"]
    }
   ]
  }`
}

func SignUrl(domain, path, key string) string {
	t := time.Now().Unix()
	raw := fmt.Sprintf("%s%d/%s", key, t, path)
	hash := utils.GetMd5(raw)
	return fmt.Sprintf("%s/%d/%s/%s", domain, t, hash, path)
}
