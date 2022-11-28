package ecode

import (
	"fmt"
	"strconv"
	"sync"
)

var (
	messages = sync.Map{}
	codes    = make(map[int]struct{})
)

// Code 错误码
// 使用接口，防止子模块通过Code(int)创建错误，从而绕过 conv 内的错误码存在检测
type Code interface {
	Code() int
	Msg() string
}

type code int

func Register(m map[Code]string) {
	for eachCode, msg := range m {
		messages.Store(eachCode.Code(), msg)
	}
}

func New(e int) Code {
	if e <= 0 {
		panic("ecode must greater than zero")
	}
	return conv(e)
}

func conv(e int) Code {
	if _, ok := codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d has already existed", e))
	}
	codes[e] = struct{}{}
	return code(e)
}

func (c code) Code() int {
	return int(c)
}

func (c code) Msg() string {
	if msg, ok := messages.Load(c.Code()); ok {
		return msg.(string)
	}
	return strconv.Itoa(c.Code())
}
