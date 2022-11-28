package middleware

import (
	"github.com/flamego/flamego"
)

type resp struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Count int64       `json:"count,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func http(r flamego.Render, code int, msg string, data interface{}, count ...int64) {
	co := int64(0)
	if len(count) > 0 {
		co = count[0]
	}
	r.JSON(code/100, &resp{
		Code:  code,
		Msg:   msg,
		Data:  data,
		Count: co,
	})
}

// HTTPSuccess 成功返回
func HTTPSuccess(r flamego.Render, data interface{}, count ...int64) {
	http(r, 20000, "success", data, count...)
}

func HTTPFail(r flamego.Render, code int, msg string, count ...int64) {
	http(r, code, msg, nil, count...)
}

func UnAuthorization(r flamego.Render) {
	HTTPFail(r, 40100, "unAuthorized")
}

func UserNotFound(r flamego.Render) {
	HTTPFail(r, 40200, "用户不存在")
}
func LoginErr(r flamego.Render) {
	HTTPFail(r, 40201, "用户名或密码错误")
}

func InValidParam(r flamego.Render) {
	HTTPFail(r, 40200, "invalid params")
}

func ServiceErr(r flamego.Render, err error) {
	HTTPFail(r, 50000, err.Error())
}
