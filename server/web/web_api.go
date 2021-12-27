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

func resultFunc(ctx iris.Context, ret interface{}) {
	if _, err := ctx.JSON(ret); err != nil {
		_, _ = ctx.Problem(newProblem(iris.StatusInternalServerError, "", err.Error()))
	}
}

func callWait(val reflect.Value, args ...reflect.Value) reflect.Value {
	f := func(val reflect.Value, args []reflect.Value) reflect.Value {
		outValue := val.Call(args)[0]
		return outValue
	}
	outValue := task.Wait(f, val, args)[0].(reflect.Value)
	return outValue
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

// 查询Token并执行
func tknCallFunc(ctx iris.Context, val reflect.Value, args ...reflect.Value) reflect.Value {
	tkn := ctx.GetHeader("token")
	if tkn == "" {
		return reflect.ValueOf(Result{Code: 2, Message: "未携带Token"})
	}

	f := func(tkn string, val reflect.Value, args []reflect.Value) reflect.Value {
		username, ok := getTknUser(tkn)
		if !ok {
			return reflect.ValueOf(Result{Code: 3, Message: "Token失效"})
		}
		nameValue := reflect.ValueOf(username)
		outValue := val.Call(append([]reflect.Value{nameValue}, args...))[0]
		return outValue
	}

	rets := task.Wait(f, tkn, val, args)
	return rets[0].(reflect.Value)
}

// 仅在队列中执行
func warpWaitHandle(fn iris.Handler) iris.Handler {
	return func(context iris.Context) {
		callWait(reflect.ValueOf(fn))
	}
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
		//log.Println(ctx.RouteName(), ctx.GetCurrentRoute().String())
		inValue, err := bodyFunc(ctx, typ.In(0))
		if err != nil {
			_, _ = ctx.Problem(newProblem(iris.StatusBadRequest, "", err.Error()))
			return
		}
		outValue := callWait(val, inValue)
		resultFunc(ctx, outValue.Interface())
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
		outValue := tknCallFunc(ctx, val)
		resultFunc(ctx, outValue.Interface())
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
		inValue, err := bodyFunc(ctx, typ.In(1))
		if err != nil {
			_, _ = ctx.Problem(newProblem(iris.StatusBadRequest, "", err.Error()))
			return
		}
		outValue := tknCallFunc(ctx, val, inValue)
		resultFunc(ctx, outValue.Interface())
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
