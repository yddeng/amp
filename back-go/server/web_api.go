package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	code     int
	result   Result
	done     chan struct{}
	doneOnce sync.Once
}

func newDone(route string) *Done {
	return &Done{
		route: route,
		code:  http.StatusOK,
		done:  make(chan struct{}),
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

func transBegin(ctx *gin.Context, fn interface{}, args ...reflect.Value) {
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
			done.code = 401
			done.result.Message = "Token验证失败"
			done.Done()
			return
		}

		ok = checkPermission(ctx, route, user)
		if !ok {
			done.code = 403
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

	ctx.JSON(done.code, done.result)
}

func getCurrentRoute(ctx *gin.Context) string {
	return ctx.FullPath()
}

func getJsonBody(ctx *gin.Context, inType reflect.Type) (inValue reflect.Value, err error) {
	if inType.Kind() == reflect.Ptr {
		inValue = reflect.New(inType.Elem())
	} else {
		inValue = reflect.New(inType)
	}
	if err = ctx.ShouldBindJSON(inValue.Interface()); err != nil {
		return
	}
	if inType.Kind() != reflect.Ptr {
		inValue = inValue.Elem()
	}
	return
}

func warpHandle(fn interface{}) gin.HandlerFunc {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	switch typ.NumIn() {
	case 2: // func(done *Done, username string)
		return func(ctx *gin.Context) {
			transBegin(ctx, fn)
		}
	case 3: // func(done *Done, username string,req struct)Result
		return func(ctx *gin.Context) {
			inValue, err := getJsonBody(ctx, typ.In(2))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Json unmarshal failed!",
					"error":   err.Error(),
				})
				return
			}

			transBegin(ctx, fn, inValue)
		}
	default:
		panic("func symbol error")
	}
}

var (
	// 允许无token的路由
	allowTokenRoute = map[string]struct{}{
		"/auth/login":      {},
		"/api/auth/login":  {},
		"/auth/logout":     {},
		"/api/auth/logout": {},
	}
	// 允许无权限的路由
	allowPermissionRoute = map[string]struct{}{
		"/auth/login":      {},
		"/api/auth/login":  {},
		"/auth/logout":     {},
		"/api/auth/logout": {},
	}
)

func checkToken(ctx *gin.Context, route string) (user string, ok bool) {
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

func checkPermission(ctx *gin.Context, route, user string) (ok bool) {
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

func initHandler(app *gin.Engine) {
	authHandle := new(authHandler)
	authGroup := app.Group("/auth")
	authGroup.POST("/login", warpHandle(authHandle.Login))
	authGroup.POST("/logout", warpHandle(authHandle.Logout))

	userHandle := new(userHandler)
	userGroup := app.Group("/user")
	userGroup.POST("/list", warpHandle(userHandle.List))
	userGroup.POST("/add", warpHandle(userHandle.Add))
	userGroup.POST("/delete", warpHandle(userHandle.Delete))

	nodeHandle := new(nodeHandler)
	nodeGroup := app.Group("/node")
	nodeGroup.POST("/list", warpHandle(nodeHandle.List))
	nodeGroup.POST("/remove", warpHandle(nodeHandle.Remove))

	cmdHandle := new(cmdHandler)
	cmdGroup := app.Group("/cmd")
	cmdGroup.POST("/list", warpHandle(cmdHandle.List))
	cmdGroup.POST("/create", warpHandle(cmdHandle.Create))
	cmdGroup.POST("/delete", warpHandle(cmdHandle.Delete))
	cmdGroup.POST("/update", warpHandle(cmdHandle.Update))
	cmdGroup.POST("/exec", warpHandle(cmdHandle.Exec))
	cmdGroup.POST("/log", warpHandle(cmdHandle.Log))

	processHandle := new(processHandler)
	processGroup := app.Group("/process")
	processGroup.POST("/glist", warpHandle(processHandle.GroupList))
	processGroup.POST("/gadd", warpHandle(processHandle.GroupAdd))
	processGroup.POST("/gremove", warpHandle(processHandle.GroupRemove))
	processGroup.POST("/list", warpHandle(processHandle.List))
	processGroup.POST("/create", warpHandle(processHandle.Create))
	processGroup.POST("/delete", warpHandle(processHandle.Delete))
	processGroup.POST("/update", warpHandle(processHandle.Update))
	processGroup.POST("/start", warpHandle(processHandle.Start))
	processGroup.POST("/stop", warpHandle(processHandle.Stop))

	kvHandle := new(kvHandler)
	kvGroup := app.Group("/kv")
	kvGroup.POST("/set", warpHandle(kvHandle.Set))
	kvGroup.POST("/get", warpHandle(kvHandle.Get))
	kvGroup.POST("/delete", warpHandle(kvHandle.Delete))

	flyfishHandle := new(flyfishHandler)
	flyfishGroup := app.Group("/flyfish")
	flyfishGroup.POST("/getMeta", warpHandle(flyfishHandle.GetMeta))
	flyfishGroup.POST("/addTable", warpHandle(flyfishHandle.AddTable))
	flyfishGroup.POST("/addField", warpHandle(flyfishHandle.AddField))
	flyfishGroup.POST("/getSetStatus", warpHandle(flyfishHandle.GetSetStatus))
	flyfishGroup.POST("/setMarkClear", warpHandle(flyfishHandle.SetMarkClear))
	flyfishGroup.POST("/addSet", warpHandle(flyfishHandle.AddSet))
	flyfishGroup.POST("/remSet", warpHandle(flyfishHandle.RemSet))
	flyfishGroup.POST("/addNode", warpHandle(flyfishHandle.AddNode))
	flyfishGroup.POST("/remNode", warpHandle(flyfishHandle.RemNode))
	flyfishGroup.POST("/addLeaderStoreToNode", warpHandle(flyfishHandle.AddLeaderStoreToNode))
	flyfishGroup.POST("/removeNodeStore", warpHandle(flyfishHandle.RemoveNodeStore))
}
