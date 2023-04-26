package hashtag

import "time_speak_server/src/exception"

type Config struct {
}

var errHashTagNotFound = exception.GqlError("标签不存在", "TAG_NOT_FOUND")
