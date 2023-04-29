package bce

import "time_speak_server/src/service/storage/utils"

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
