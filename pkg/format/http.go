package format

import (
	"github.com/flamego/flamego"
	"supreme-flamego/pkg/ecode"
	"net/http"
)

type resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// HTTP code 为错误码，可以传入公共错误码，也可以传入模块内部自定义错误码
func HTTP(r flamego.Render, code ecode.Code, data interface{}) {
	status := http.StatusOK
	if code.Code() >= 10000 {
		status = code.Code() / 100
	}
	r.JSON(status, &resp{
		Code: code.Code(),
		Msg:  code.Msg(),
		Data: data,
	})
}

// HTTPSuccess 成功返回
func HTTPSuccess(r flamego.Render, data interface{}) {
	HTTP(r, ecode.OK, data)
}

// HTTPInvalidParams 参数错误或存在非法字符
func HTTPInvalidParams(r flamego.Render) {
	HTTP(r, ecode.InvalidParams, nil)
}

func HTTPServerErr(r flamego.Render) {
	HTTP(r, ecode.ServerErr, nil)
}

func HTTPUnauthorized(r flamego.Render) {
	HTTP(r, ecode.Unauthorized, nil)
}
