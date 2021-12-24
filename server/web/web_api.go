package web

import (
	"github.com/kataras/iris/v12"
	"initial-sever/task"
	"reflect"
)

// 应答结构
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func callResult(ctx iris.Context, ret interface{}) {
	if _, err := ctx.JSON(ret); err != nil {
		_, _ = ctx.Problem(newProblem(iris.StatusInternalServerError, "", err.Error()))
	}
}

func bodyFunc(ctx iris.Context, inType reflect.Type) (inValue reflect.Value, err error) {
	if inType.Kind() == reflect.Ptr {
		inValue = reflect.New(inType.Elem())
	} else {
		inValue = reflect.New(inType)
	}
	if err = ctx.ReadJSON(inValue.Interface()); err != nil {
		return
	}
	if inType.Kind() != reflect.Ptr {
		inValue = inValue.Elem()
	}
	return
}

func tokenFunc(ctx iris.Context) (reflect.Value, *Result) {
	var nameValue reflect.Value
	tkn := ctx.GetHeader("token")
	if tkn == "" {
		return nameValue, &Result{Code: 2, Message: "未携带Token"}
	}

	rets := task.Wait(getToken, tkn)
	username, ok := rets[0].(string), rets[1].(bool)
	if !ok {
		return nameValue, &Result{Code: 3, Message: "Token失效"}
	}
	nameValue = reflect.ValueOf(username)
	return nameValue, nil
}

// func(req struct)Result // body解析对应的json
func warpJsonHandle(fn interface{}) iris.Handler {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != 1 || typ.NumOut() != 1 {
		panic("func symbol error")
	}
	return func(ctx iris.Context) {
		inValue, err := bodyFunc(ctx, typ.In(0))
		if err != nil {
			_, _ = ctx.Problem(newProblem(iris.StatusBadRequest, "", err.Error()))
			return
		}
		outValue := val.Call([]reflect.Value{inValue})[0]
		callResult(ctx, outValue.Interface())
	}
}

// func(username string)Result // 仅验证token
func warpTokenHandle(fn interface{}) iris.Handler {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != 1 || typ.NumOut() != 1 {
		panic("func symbol error")
	}
	return func(ctx iris.Context) {
		nameValue, ret := tokenFunc(ctx)
		if ret != nil {
			callResult(ctx, ret)
			return
		}
		outValue := val.Call([]reflect.Value{nameValue})[0]
		callResult(ctx, outValue.Interface())
	}
}

// func(username string,req struct)Result //  body解析对应的json, 验证token
func warpTokenJsonHandle(fn interface{}) iris.Handler {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != 2 || typ.NumOut() != 1 {
		panic("func symbol error")
	}
	return func(ctx iris.Context) {
		nameValue, ret := tokenFunc(ctx)
		if ret != nil {
			callResult(ctx, ret)
			return
		}
		inValue, err := bodyFunc(ctx, typ.In(1))
		if err != nil {
			_, _ = ctx.Problem(newProblem(iris.StatusBadRequest, "", err.Error()))
			return
		}
		outValue := val.Call([]reflect.Value{nameValue, inValue})[0]
		callResult(ctx, outValue.Interface())
	}
}

func newProblem(statusCode int, title, detail string) iris.Problem {
	p := iris.NewProblem().Status(statusCode)
	if title != "" {
		p.Title(title)
	}
	p.Detail(detail)
	return p
}

func handleCORS(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Next()
}
