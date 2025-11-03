package errs

//用户模块

const (
	UserInternalServerError = 501001
	UserInvalidInput        = 401001
	// UserInvalidOrPassword 用户不存在或者密码错误，这个你要小心，
	// 防止有人跟你过不去
	UserInvalidOrPassword = 401002
	UserDuplicateEmail    = 401003
	UserNotFound          = 401004
	UserCodeSendTooMany   = 401005
	UserVerifyCodeErr     = 401006
)
