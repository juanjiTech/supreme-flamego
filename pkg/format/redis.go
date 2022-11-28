package format

import (
	"bytes"
	"strings"
	"sync"
)

var keyPool = sync.Pool{
	New: func() interface{} {
		b := &bytes.Buffer{}
		b.Grow(32)
		return b
	},
}

type key struct{}

var Key key

func (key) Gen(business string, indexes ...string) string {
	buf := keyPool.Get().(*bytes.Buffer)

	buf.WriteString(business)
	buf.WriteString(":")
	buf.WriteString(strings.Join(indexes, ":"))

	key := buf.String()
	buf.Reset()
	keyPool.Put(buf)
	return key
}

func (key) AuthToken(token string) string {
	return Key.Gen("au", "t", token)
}

func (key) AuthCallback(mark string) string {
	return Key.Gen("au", "c", mark)
}

func (key) QRLogin(mark string) string {
	return Key.Gen("au", "qr", mark)
}

func (key) ServiceTicket(mark string) string {
	return Key.Gen("cas", "st", mark)
}
