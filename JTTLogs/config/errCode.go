package config

const (
	AdminUser = "admin" //管理员
	AdminPwd  = "admin" //管理员密码
	GuestUser = "guest" //游客用户
	GuestPwd  = "guest" //游客密码
)

const (
	JTT_ERROR_SUCCESS_OK        = 200
	JTT_ERROR_FAILED            = 400
	JTT_ERROR_PARAMETER_ERROR   = 300
	JTT_ERROR_INVALID_TOKEN     = 650
	JTT_ERROR_SERVER_ERROR      = 402
	JTT_ERROR_OVERTIME          = 403
	JTT_ERROR_DEVICELOST        = 406
	JTT_ERROR_USER_NOT_FOUND    = 600
	JTT_ERROR_CHANNEL_NOT_FOUNT = 602
	JTT_ERROR_URL_NOT_FOUND     = 603
	JTT_ERROR_DB_ERR            = 401
)

var ErrCodeMap map[int64]string

func init() {
	ErrCodeMap = map[int64]string{
		JTT_ERROR_SUCCESS_OK:        "Success OK",
		JTT_ERROR_FAILED:            "Failed",
		JTT_ERROR_PARAMETER_ERROR:   "Parameter error",
		JTT_ERROR_INVALID_TOKEN:     "Invalid Session", //无效的session
		JTT_ERROR_SERVER_ERROR:      "Server Error",
		JTT_ERROR_OVERTIME:          "Request over time",
		JTT_ERROR_DEVICELOST:        "Device Not Online",
		JTT_ERROR_USER_NOT_FOUND:    "User or password error",
		JTT_ERROR_URL_NOT_FOUND:     "Cannot get url response",
		JTT_ERROR_DB_ERR:            "DB err",
		JTT_ERROR_CHANNEL_NOT_FOUNT: "Channel not fount",
	}
}
