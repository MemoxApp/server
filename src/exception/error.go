package exception

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
	"time_speak_server/src/log"
)

var (
	ErrUsernameInvalid  = GqlError("用户名不合法", "USERNAME_INVALID")
	ErrUsernameOccupied = GqlError("用户名已被使用", "USERNAME_OCCUPIED")
	ErrMailOccupied     = GqlError("邮箱已被使用", "MAIL_OCCUPIED")

	ErrVerifyCodeWrong         = GqlError("验证码错误", "VERIFY_CODE_WRONG")
	ErrUsernameOrPasswordWrong = GqlError("用户名或密码错误", "USERNAME_OR_PASSWORD_WRONG")
	ErrUserNotFound            = GqlError("用户不存在", "USER_NOT_FOUND")
	ErrReplyNotFound           = GqlError("回复不存在", "REPLY_NOT_FOUND")

	ErrTooManyRequest = GqlError("请求过于频繁", "TOO_MANY_REQUEST")
	ErrParamInvalid   = GqlError("参数不合法", "PARAM_INVALID")
	ErrContentEmpty   = GqlError("你还什么都没有写呢", "CONTENT_EMPTY")
	ErrInvalidID      = GqlError("错误的ID", "INVALID_ID")

	ErrTitleTooLong   = GqlError("标题太长啦", "TITLE_TOO_LONG")
	ErrContentTooLong = GqlError("内容太长啦", "CONTENT_TOO_LONG")
	ErrContentExist   = GqlError("内容已存在", "CONTENT_EXIST")
	ErrSubscribeExist = GqlError("订阅已存在", "SUBSCRIBE_EXIST")

	ErrPermissionDenied = GqlError("访问权限不足", "PERMISSION_DENIED")
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
		Extensions: map[string]interface{}{
			"code": code,
		},
	}
}
