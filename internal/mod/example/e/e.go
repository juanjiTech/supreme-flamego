// Package e 主要用于定义app特有的错误码
package e

import "supreme-flamego/pkg/ecode"

var (
	ErrNotExist = ecode.New(40400)
)

var ECode = map[ecode.Code]string{
	ErrNotExist: "资源不存在",
}

func init() {
	ecode.Register(ECode) // 使用错误码注册机制保证软件全局不存在错误码冲突
}
