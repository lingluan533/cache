package config

// 业务状态码
const (
	EdgeNodeNotFound    = "10000"
	UserRepeatRegister  = "50001"
	UserRegisterSuccess = "50002"
	ServerInternalError = "50000"
	PasswordNotSame     = "50003"
	UserLoginSuccess    = "60002"
)

// 错误提示信息
const (
	MsgUserRepeatRegister  = "用户名已存在！"
	MsgUserRegisterSuccess = "用户注册成功！"
	MsgServerInternalError = "服务器内部错误，请联系管理员！"
	MsgPasswordNotSame     = "两次输入密码不一致，请重试!"
	MsgUserLoginSuccess    = "用户登录成功！"
)
