// Package ecode 公共错误码
package ecode

var (
	OK               = conv(20000)
	UnknownErr       = conv(50000)
	ServerErr        = conv(50001)
	InvalidParams    = conv(40000)
	RateLimit        = conv(42900)
	Unauthorized     = conv(40100)
	PermissionDenied = conv(40301)
)

func init() {
	Register(map[Code]string{
		OK:               "ok",
		UnknownErr:       "未知错误",
		ServerErr:        "服务器错误",
		InvalidParams:    "参数错误或存在非法字符",
		RateLimit:        "访问频率过高",
		Unauthorized:     "身份验证失败",
		PermissionDenied: "权限不足",
	})
}
