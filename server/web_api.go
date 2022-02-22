package server

import (
	"github.com/kataras/iris/v12"
	"reflect"
	"sync"
)

// 应答结构
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

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
		if this.result.Message != "" {
			this.result.Code = 1
		}
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
	if err := taskQueue.SubmitTask(webTask(func() {
		user, ok := checkToken(ctx, route)
		if !ok {
			done.statue = 401
			done.result.Message = "Token验证失败"
			done.Done()
			return
		}

		ok = checkPermission(ctx, route, user)
		if !ok {
			done.statue = 403
			done.result.Message = "无操作权限"
			done.Done()
			return
		}
		val.Call(append([]reflect.Value{reflect.ValueOf(done), reflect.ValueOf(user)}, args...))
	}), true); err != nil {
		done.result.Message = "当前访问人数过多"
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
	}
)

func handleCORS(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Next()
}

func checkToken(ctx iris.Context, route string) (user string, ok bool) {
	if _, ok = allowTokenRoute[route]; ok {
		return
	}
	tkn := ctx.GetHeader("Access-Token")
	if tkn == "" {
		ok = false
		return
	}
	if user, ok = getTknUser(tkn); !ok {
		ok = false
		return
	}
	ok = true
	return
}

func checkPermission(ctx iris.Context, route, user string) (ok bool) {
	if _, ok = allowPermissionRoute[route]; ok {
		return
	}
	ok = true
	return
}

func listRange(pageNo, pageSize, length int) (start int, end int) {
	start = (pageNo - 1) * pageSize
	if start < 0 {
		start = 0
	}
	if start > length {
		start = length
	}
	end = start + pageSize
	if end > length {
		end = length
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
	userRouter.Post("/list", warpHandle(userHandle.List))
	userRouter.Post("/add", warpHandle(userHandle.Add))
	userRouter.Post("/delete", warpHandle(userHandle.Delete))

	nodeHandle := new(nodeHandler)
	nodeRouter := app.Party("/node")
	nodeRouter.Post("/list", warpHandle(nodeHandle.List))
	nodeRouter.Post("/remove", warpHandle(nodeHandle.Remove))

	cmdHandle := new(cmdHandler)
	cmdRouter := app.Party("/cmd")
	cmdRouter.Post("/list", warpHandle(cmdHandle.List))
	cmdRouter.Post("/create", warpHandle(cmdHandle.Create))
	cmdRouter.Post("/delete", warpHandle(cmdHandle.Delete))
	cmdRouter.Post("/update", warpHandle(cmdHandle.Update))
	cmdRouter.Post("/exec", warpHandle(cmdHandle.Exec))
	cmdRouter.Post("/log", warpHandle(cmdHandle.Log))

	processHandle := new(processHandler)
	processRouter := app.Party("/process")
	processRouter.Post("/glist", warpHandle(processHandle.GroupList))
	processRouter.Post("/gadd", warpHandle(processHandle.GroupAdd))
	processRouter.Post("/gremove", warpHandle(processHandle.GroupRemove))
	processRouter.Post("/list", warpHandle(processHandle.List))
	processRouter.Post("/create", warpHandle(processHandle.Create))
	processRouter.Post("/delete", warpHandle(processHandle.Delete))
	processRouter.Post("/update", warpHandle(processHandle.Update))
	processRouter.Post("/start", warpHandle(processHandle.Start))
	processRouter.Post("/stop", warpHandle(processHandle.Stop))

	kvHandle := new(kvHandler)
	kvRouter := app.Party("/kv")
	kvRouter.Post("/set", warpHandle(kvHandle.Set))
	kvRouter.Post("/get", warpHandle(kvHandle.Get))
	kvRouter.Post("/delete", warpHandle(kvHandle.Delete))

	flyfishHandle := new(flyfishHandler)
	flyfishRouter := app.Party("/flyfish")
	flyfishRouter.Post("/getMeta", warpHandle(flyfishHandle.GetMeta))
	flyfishRouter.Post("/addTable", warpHandle(flyfishHandle.AddTable))
	flyfishRouter.Post("/addField", warpHandle(flyfishHandle.AddField))
	flyfishRouter.Post("/getSetStatus", warpHandle(flyfishHandle.GetSetStatus))
	flyfishRouter.Post("/setMarkClear", warpHandle(flyfishHandle.SetMarkClear))
	flyfishRouter.Post("/addSet", warpHandle(flyfishHandle.AddSet))
	flyfishRouter.Post("/remSet", warpHandle(flyfishHandle.RemSet))
	flyfishRouter.Post("/addNode", warpHandle(flyfishHandle.AddNode))
	flyfishRouter.Post("/remNode", warpHandle(flyfishHandle.RemNode))
	flyfishRouter.Post("/addLeaderStoreToNode", warpHandle(flyfishHandle.AddLeaderStoreToNode))
	flyfishRouter.Post("/removeNodeStore", warpHandle(flyfishHandle.RemoveNodeStore))
}
