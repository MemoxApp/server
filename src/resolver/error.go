package resolver

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
	"time_speak_server/src/log"
)

var (
	errUsernameInvalid  = gqlError("用户名不合法", "USERNAME_INVALID")
	errUsernameOccupied = gqlError("用户名已被使用", "USERNAME_OCCUPIED")
	errMailOccupied     = gqlError("邮箱已被使用", "MAIL_OCCUPIED")

	errVerifyCodeWrong         = gqlError("验证码错误", "VERIFY_CODE_WRONG")
	errUsernameOrPasswordWrong = gqlError("用户名或密码错误", "USERNAME_OR_PASSWORD_WRONG")

	errUserNotFound  = gqlError("用户不存在", "USER_NOT_FOUND")
	errPostNotFound  = gqlError("帖子不存在", "POST_NOT_FOUND")
	errReplyNotFound = gqlError("回复不存在", "REPLY_NOT_FOUND")

	errTooManyRequest = gqlError("请求过于频繁", "TOO_MANY_REQUEST")
	errParamInvalid   = gqlError("参数不合法", "PARAM_INVALID")
	errContentEmpty   = gqlError("你还什么都没有写呢", "CONTENT_EMPTY")

	errTitleTooLong   = gqlError("标题太长啦", "TITLE_TOO_LONG")
	errContentTooLong = gqlError("内容太长啦", "Content_TOO_LONG")

	ErrPermissionDenied = gqlError("访问权限不足", "PERMISSION_DENIED")
)

func internalError(err error) error {
	log.Error("internalError", "error", err)
	return &gqlerror.Error{
		Message: "服务器开小差了，等会再试试吧",
		Extensions: map[string]interface{}{
			"code":  "INTERNAL_ERROR",
			"debug": err,
		},
	}
}

func gqlError(message, code string) error {
	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			"code": code,
		},
	}
}
