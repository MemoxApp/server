package bce

import (
	"fmt"
	"memox_server/src/service/storage/local"
	"memox_server/src/service/storage/utils"
	"time"
)

func (b *BCE) getWritePermissionACL(userID, filename, bucketName string) string {
	return ` {
  "accessControlList": [
   {   
    "service":"bce:bos",
    "region":"` + b.Config.Region + `",
    "effect": "Allow",
    "resource": ["` + bucketName + "/" + utils.GenerateResourcePath(userID, filename) + `"],
    "permission": ["WRITE"]
    }
   ]
  }`
}

func SignUrlTypeA(domain, fileName, uid, key string) string {
	rand := local.GenerateRandomString(6)
	timestamp := time.Now().Unix()
	raw := fmt.Sprintf("/%s-%d-%s-%s-%s", fileName, timestamp, rand, uid, key)
	hash := utils.GetMd5(raw)
	return fmt.Sprintf("%s/%s?auth_key=%d-%s-%s-%s", domain, fileName, timestamp, rand, uid, hash)
}

func SignUrlTypeB(domain, fileName, key string) string {
	timestamp := time.Now().Unix()
	raw := fmt.Sprintf("%s%d/%s", key, timestamp, fileName)
	hash := utils.GetMd5(raw)
	return fmt.Sprintf("%s/%d/%s/%s", domain, timestamp, hash, fileName)
}

func SignUrlTypeC(domain, fileName, key string) string {
	timestamp := time.Now().Unix()
	raw := fmt.Sprintf("%s/%s%d", key, fileName, timestamp)
	hash := utils.GetMd5(raw)
	return fmt.Sprintf("%s/%s?md5hash=%s&timestamp=%d", domain, fileName, hash, timestamp)
}
