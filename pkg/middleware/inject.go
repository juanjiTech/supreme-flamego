package middleware

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/flamego/flamego"
	"reflect"
	"supreme-flamego/internal/cache"
	"supreme-flamego/internal/database"
	"supreme-flamego/pkg/format"
	"supreme-flamego/pkg/utils/page"
	"supreme-flamego/pkg/websocket"
)

func InjectDB(key ...string) flamego.Handler {
	if len(key) > 1 {
		panic("InjectDB only accept at most one key")
	}
	if len(key) == 1 {
		return func(c flamego.Context) {
			c.Map(database.GetDB(key[0]))
		}
	}
	return func(c flamego.Context) {
		c.Map(database.GetDB("*"))
	}
}

func InjectWebsocket(key ...string) flamego.Handler {
	if len(key) > 1 {
		panic("InjectWebsocket only accept at most one key")
	}
	if len(key) == 1 {
		return func(c flamego.Context) {
			c.Map(websocket.GetSocketManager(key[0]))
		}
	}
	return func(c flamego.Context) {
		c.Map(websocket.GetSocketManager("*"))
	}
}

func InjectCache(key ...string) flamego.Handler {
	if len(key) > 1 {
		panic("InjectCache only accept at most one key")
	}
	if len(key) == 1 {
		return func(c flamego.Context) {
			c.Map(cache.GetCache(key[0]))
		}
	}
	return func(c flamego.Context) {
		c.Map(cache.GetCache())
	}
}

func InjectJson[T any]() flamego.Handler {
	return func(r flamego.Render, c flamego.Context) {
		var req T
		body, err := c.Request().Body().Bytes()
		if err = json.Unmarshal(body, &req); err != nil {
			format.HTTPInvalidParams(r)
			return
		}
		_, err = govalidator.ValidateStruct(&req)
		if err != nil {
			format.HTTPInvalidParams(r)
			return
		}
		c.Map(req)
	}
}

func InjectQuery[T any]() flamego.Handler {
	return func(r flamego.Render, c flamego.Context) {
		var req T
		t := reflect.TypeOf(req)
		v := reflect.ValueOf(&req).Elem()
		var tag reflect.StructTag
		for i := 0; i < v.NumField(); i++ {
			tag = t.Field(i).Tag
			if value, ok := tag.Lookup("query"); ok && t.Field(i).IsExported() {
				switch v.Field(i).Kind() {
				case reflect.String:
					v.Field(i).SetString(c.Query(value))
				case reflect.Int, reflect.Int64:
					v.Field(i).SetInt(c.QueryInt64(value))
				case reflect.Bool:
					v.Field(i).SetBool(c.QueryBool(value))
				}
			}
		}

		_, err := govalidator.ValidateStruct(&req)
		if err != nil {
			format.HTTPInvalidParams(r)
			return
		}
		c.Map(req)
	}
}

func InjectParam[T any]() flamego.Handler {
	return func(r flamego.Render, c flamego.Context) {
		var req T
		t := reflect.TypeOf(req)
		v := reflect.ValueOf(&req).Elem()
		var tag reflect.StructTag
		for i := 0; i < v.NumField(); i++ {
			tag = t.Field(i).Tag
			if value, ok := tag.Lookup("param"); ok && t.Field(i).IsExported() {
				switch v.Field(i).Kind() {
				case reflect.String:
					v.Field(i).SetString(c.Param(value))
				case reflect.Int, reflect.Int64:
					v.Field(i).SetInt(c.ParamInt64(value))
				case reflect.Bool:
					panic("bool type not support")
				}
			}
		}

		_, err := govalidator.ValidateStruct(&req)
		if err != nil {
			format.HTTPInvalidParams(r)
			return
		}
		c.Map(req)
	}
}

func InjectPaginate() flamego.Handler {
	return func(r flamego.Render, c flamego.Context) {
		var req page.Paginate
		req.Current = c.QueryInt("current")
		req.PageSize = c.QueryInt("pageSize")
		c.Map(req)
	}
}
