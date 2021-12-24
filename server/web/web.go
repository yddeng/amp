package web

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"log"
	"reflect"
	"time"
)

type Config struct {
	WebAddress string `json:"web_address"`
	FilePath   string `json:"file_path"`
}

func handleCORS(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")
	ctx.Next()
}

func HandleJSONJSON(fn interface{}) iris.Handler {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != 1 || typ.NumOut() != 1 {
		panic("func symbol error")
	}

	return func(ctx iris.Context) {
		inType := typ.In(0)
		var inValue reflect.Value
		if inType.Kind() == reflect.Ptr {
			inValue = reflect.New(inType.Elem())
		} else {
			inValue = reflect.New(inType)
		}
		if err := ctx.ReadJSON(inValue.Interface()); err != nil {
			ctx.Problem(NewProblem(iris.StatusBadRequest, "", err.Error()))
			return
		}

		if inType.Kind() != reflect.Ptr {
			inValue = inValue.Elem()
		}

		outValue := val.Call([]reflect.Value{inValue})[0]
		if _, err := ctx.JSON(outValue.Interface()); err != nil {
			ctx.Problem(NewProblem(iris.StatusInternalServerError, "", err.Error()))
		}
	}
}

func NewProblem(statusCode int, title, detail string) iris.Problem {
	p := iris.NewProblem().Status(statusCode)
	if title != "" {
		p.Title(title)
	}
	p.Detail(detail)
	return p
}

type token struct {
	user   string
	expire time.Time
}

var (
	tknUser = map[string]*token{}
	dur     = time.Hour
)

func addToken(u string) (tkn string) {
	now := time.Now()
	tkn = fmt.Sprintf("%d", now.UnixNano())
	t := &token{
		user:   u,
		expire: now.Add(dur),
	}
	tknUser[tkn] = t
	return
}

func getToken(tkn string) (string, bool) {
	t, ok := tknUser[tkn]
	if !ok {
		return "", false
	}
	t.expire = time.Now().Add(dur)
	return t.user, true
}

func CheckToken(ctx iris.Context) (string, *Result) {
	tkn := ctx.GetHeader("token")
	if tkn == "" {
		return "", &Result{
			Code:    2,
			Message: "未携带Token",
		}
	}

	username, ok := getToken(tkn)
	if !ok {
		return "", &Result{
			Code:    3,
			Message: "Token失效",
		}
	}

	return username, nil
}

// 应答结构
type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Handler struct{}

func RunWeb(cfg *Config) {
	app := iris.New()
	app.Use(logger.New())
	// 跨域
	app.Use(handleCORS)

	handler := new(Handler)
	authRouter := app.Party("/auth")
	authRouter.Post("/login", handler.Login)

	userRouter := app.Party("/user")
	userRouter.Post("/getNav", handler.GetNav)

	log.Printf("web server run %s.\n", cfg.WebAddress)
	if err := app.Listen(cfg.WebAddress); err != nil {
		panic(err)
	}
}

func init() {

}
