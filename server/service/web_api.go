package service

import (
	"github.com/kataras/iris/v12"
	"github.com/yddeng/utils/task"
	"reflect"
	"sync"
)

// 应答结构
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var webTransQueue = task.NewTaskPool(1, 2048)

type Done struct {
	route    string
	statue   int
	result   Result
	done     chan struct{}
	doneOnce sync.Once
}

func newDone(route string) *Done {
	return &Done{
		route:  route,
		statue: 200,
		done:   make(chan struct{}),
	}
}

func (this *Done) Done() {
	this.doneOnce.Do(func() {
		close(this.done)
	})
}

func (this *Done) Wait() {
	<-this.done
}

type webTask func()

func (t webTask) Do() {
	t()
}

func transBegin(ctx iris.Context, fn interface{}, args ...reflect.Value) {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != len(args)+2 {
		panic("func argument error")
	}

	route := getCurrentRoute(ctx)
	done := newDone(route)
	if err := webTransQueue.SubmitTask(webTask(func() {
		user, ret := checkToken(ctx, route)
		if ret.Code != 0 {
			done.statue = 401
			done.result.Code = ret.Code
			done.result.Message = ret.Message
			done.Done()
			return
		}

		ret = checkPermission(ctx, route, user)
		if ret.Code != 0 {
			done.statue = 403
			done.result.Code = ret.Code
			done.result.Message = ret.Message
			done.Done()
			return
		}
		val.Call(append([]reflect.Value{reflect.ValueOf(done), reflect.ValueOf(user)}, args...))
	}), true); err != nil {
		done.result.Code = 1
		done.result.Message = "当前访问人数"
		done.Done()
	}
	done.Wait()

	if done.statue != 200 {
		ctx.StatusCode(done.statue)
	}
	if _, err := ctx.JSON(done.result); err != nil {
		_, _ = ctx.Problem(newProblem(iris.StatusInternalServerError, "", err.Error()))
	}
}

func getCurrentRoute(ctx iris.Context) string {
	return ctx.GetCurrentRoute().Path()
}

func getJsonBody(ctx iris.Context, inType reflect.Type) (inValue reflect.Value, err error) {
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

func warpHandle(fn interface{}) iris.Handler {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	switch typ.NumIn() {
	case 2: // func(done *Done, username string)
		return func(ctx iris.Context) {
			transBegin(ctx, fn)
		}
	case 3: // func(done *Done, username string,req struct)Result
		return func(ctx iris.Context) {
			inValue, err := getJsonBody(ctx, typ.In(2))
			if err != nil {
				_, _ = ctx.Problem(newProblem(iris.StatusBadRequest, "", err.Error()))
				return
			}

			transBegin(ctx, fn, inValue)
		}
	default:
		panic("func symbol error")
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

var (
	// 允许无token的路由
	allowTokenRoute = map[string]struct{}{
		"/auth/login":  {},
		"/auth/logout": {},
	}
	// 允许无权限的路由
	allowPermissionRoute = map[string]struct{}{
		"/auth/login":  {},
		"/auth/logout": {},
		"/user/nav":    {},
		"/user/info":   {},
	}
)

func handleCORS(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Next()
}

func checkToken(ctx iris.Context, route string) (user string, ret Result) {
	var ok bool
	if _, ok = allowTokenRoute[route]; ok {
		return
	}
	tkn := ctx.GetHeader("Access-Token")
	if tkn == "" {
		ret.Code = 401
		ret.Message = "未携带Token"
		return
	}
	if user, ok = getTknUser(tkn); !ok {
		ret.Code = 401
		ret.Message = "Token失效"
		return
	}
	return
}

func checkPermission(ctx iris.Context, route, user string) (ret Result) {
	var ok bool
	if _, ok = allowPermissionRoute[route]; ok {
		return
	}
	return
}

func initHandler(app *iris.Application) {
	authHandle := new(authHandler)
	authRouter := app.Party("/auth")
	authRouter.Post("/login", warpHandle(authHandle.Login))
	authRouter.Post("/logout", warpHandle(authHandle.Logout))

	userHandle := new(userHandler)
	userRouter := app.Party("/user")
	userRouter.Get("/nav", warpHandle(userHandle.Nav))
	userRouter.Get("/info", warpHandle(userHandle.Info))
	userRouter.Get("/list", warpHandle(userHandle.List))
	userRouter.Post("/add", warpHandle(userHandle.Add))
	userRouter.Post("/delete", warpHandle(userHandle.Delete))

	nodeHandle := new(nodeHandler)
	nodeRouter := app.Party("/node")
	nodeRouter.Get("/list", warpHandle(nodeHandle.List))

	{
		projectRouter := app.Party("/project")

		clusterHandle := new(clusterHandler)
		clusterRouter := projectRouter.Party("/cluster")
		clusterRouter.Get("/list", warpHandle(clusterHandle.List))
		clusterRouter.Post("/create", warpHandle(clusterHandle.Create))
		clusterRouter.Post("/delete", warpHandle(clusterHandle.Delete))

		itemHandle := new(itemHandler)
		itemRoute := clusterRouter.Party("/item")
		itemRoute.Get("/list", warpHandle(itemHandle.List))
		itemRoute.Post("/create", warpHandle(itemHandle.Create))
		itemRoute.Post("/delete", warpHandle(itemHandle.Delete))
		itemRoute.Post("/start", warpHandle(itemHandle.Start))
		itemRoute.Post("/single", warpHandle(itemHandle.Single))

		templateHandle := new(templateHandler)
		templateRoute := projectRouter.Party("/template")
		templateRoute.Get("/list", warpHandle(templateHandle.List))
		templateRoute.Post("/create", warpHandle(templateHandle.Create))
		templateRoute.Post("/delete", warpHandle(templateHandle.Delete))
		templateRoute.Post("/update", warpHandle(templateHandle.Update))

	}
}
