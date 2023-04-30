package exception

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
	"time_speak_server/src/log"
)

var (
	ErrUsernameInvalid        = GqlError("用户名不合法", "USERNAME_INVALID")
	ErrMailOccupied           = GqlError("邮箱已被使用", "MAIL_OCCUPIED")
	ErrVerifyCodeWrong        = GqlError("验证码错误", "VERIFY_CODE_WRONG")
	ErrEmailOrPasswordWrong   = GqlError("邮箱或密码错误", "USERNAME_OR_PASSWORD_WRONG")
	ErrMailNotFound           = GqlError("邮箱不存在", "MAIL_NOT_FOUND")
	ErrUserNotFound           = GqlError("用户不存在", "USER_NOT_FOUND")
	ErrReplyNotFound          = GqlError("回复不存在", "REPLY_NOT_FOUND")
	ErrResourceNotFound       = GqlError("资源不存在", "RESOURCE_NOT_FOUND")
	ErrTooManyRequest         = GqlError("请求过于频繁", "TOO_MANY_REQUEST")
	ErrContentEmpty           = GqlError("你还什么都没有写呢", "CONTENT_EMPTY")
	ErrInvalidID              = GqlError("错误的ID", "INVALID_ID")
	ErrInvalidFileName        = GqlError("不合法的文件名", "INVALID_FILE_NAME")
	ErrDeleteResource         = GqlError("删除资源文件错误", "DELETE_RESOURCE_ERROR")
	ErrInvalidStorageProvider = GqlError("不合法的存储提供者", "INVALID_STORAGE_PROVIDER")
	ErrMemoryAlreadyArchived  = GqlError("记忆已归档", "MEMORY_ALREADY_ARCHIVED")
	ErrTitleTooLong           = GqlError("标题太长啦", "TITLE_TOO_LONG")
	ErrContentTooLong         = GqlError("内容太长啦", "CONTENT_TOO_LONG")
	ErrContentExist           = GqlError("内容已存在", "CONTENT_EXIST")
	ErrResourceExist          = GqlError("资源已存在", "RESOURCE_EXIST")
	ErrSubscribeExist         = GqlError("订阅已存在", "SUBSCRIBE_EXIST")
	ErrInvalidSession         = GqlError("无效的Session", "INVALID_SESSION")
	ErrCommentNotFound        = GqlError("回复不存在", "COMMENT_NOT_FOUND")
	ErrMemoryNotFound         = GqlError("记忆不存在", "MEMORY_NOT_FOUND")
	ErrCommentNotArchived     = GqlError("回复未归档，无法删除", "COMMENT_NOT_ARCHIVED")
	ErrPermissionDenied       = GqlError("访问权限不足", "PERMISSION_DENIED")
	ErrHashTagNotFound        = GqlError("话题不存在", "HASH_TAG_NOT_FOUND")
	ErrHashTagNotArchived     = GqlError("话题未归档，无法删除", "HASH_TAG_NOT_ARCHIVED")
	ErrHashTagHasMemories     = GqlError("话题下存在记忆，无法删除", "HASH_TAG_HAS_MEMORIES")
	ErrMemoryNotArchived      = GqlError("记忆未归档，无法删除", "MEMORY_NOT_ARCHIVED")
	ErrResourceHasReference   = GqlError("资源被引用，无法删除", "RESOURCE_HAS_REFERENCE")
	ErrCommentArchived        = GqlError("回复已归档，无法评论", "COMMENT_ARCHIVED")
)

func InternalError(err error) error {
	log.Error("internalError", "error", err)
	return &gqlerror.Error{
		Message: "服务器开小差了，等会再试试吧",
		Extensions: map[string]interface{}{
			"code":  "INTERNAL_ERROR",
			"debug": err,
		},
	}
}

func GqlError(message, code string) error {
	return &gqlerror.Error{
		Message: message,
	}
}
