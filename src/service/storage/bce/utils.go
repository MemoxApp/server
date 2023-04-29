package bce

func (b *BCE) getWritePermissionACL(userID, filename string) string {
	return ` {
  "accessControlList": [
   {   
    "service":"bce:bos",
    "region":"` + b.Config.Region + `",
    "effect": "Allow",
    "resource": ["users/` + userID + "/" + filename + `"],
    "permission": ["WRITE"]
    }
   ]
  }`
}
