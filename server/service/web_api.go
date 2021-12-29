package service

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
func tknCallFunc(ctx iris.Context, val reflect.Value, args ...reflect.Value) (reflect.Value, int) {
	tkn := ctx.GetHeader("Access-Token")
	if tkn == "" {
		return reflect.ValueOf(Result{Code: 1, Message: "未携带Token"}), 401
	}

	f := func(tkn string, val reflect.Value, args []reflect.Value) (reflect.Value, int) {
		username, ok := getTknUser(tkn)
		if !ok {
			return reflect.ValueOf(Result{Code: 1, Message: "Token失效"}), 401
		}
		nameValue := reflect.ValueOf(username)
		outValue := val.Call(append([]reflect.Value{nameValue}, args...))[0]
		return outValue, 200
	}

	rets := task.Wait(f, tkn, val, args)
	return rets[0].(reflect.Value), rets[1].(int)
}

// 仅在队列中执行
func warpWaitHandle(fn func(ctx iris.Context) Result) iris.Handler {
	return func(ctx iris.Context) {
		outValue := callWait(reflect.ValueOf(fn), reflect.ValueOf(ctx))
		resultFunc(ctx, outValue.Interface())
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
		outValue, statue := tknCallFunc(ctx, val)
		ctx.StatusCode(statue)
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

		outValue, statue := tknCallFunc(ctx, val, inValue)
		ctx.StatusCode(statue)
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

func initHandler(app *iris.Application) {
	authHandle := new(Auth)
	authRouter := app.Party("/auth")
	authRouter.Post("/login", warpJsonHandle(authHandle.Login))
	authRouter.Post("/logout", warpWaitHandle(authHandle.Logout))

	userHandle := new(User)
	userRouter := app.Party("/user")
	userRouter.Get("/nav", warpTokenHandle(userHandle.Nav))
	userRouter.Get("/info", warpTokenHandle(userHandle.Info))
	userRouter.Post("/list", warpTokenJsonHandle(userHandle.List))
	userRouter.Post("/add", warpTokenJsonHandle(userHandle.Add))
	userRouter.Post("/delete", warpTokenJsonHandle(userHandle.Delete))

	nodeHandle := new(Node)
	nodeRouter := app.Party("/node")
	nodeRouter.Get("/list", warpTokenHandle(nodeHandle.List))
}
