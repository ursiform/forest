package forest

import (
	"fmt"
	"github.com/ursiform/bear"
	"log"
	"net/http"
	"time"
)

type App struct {
	Debug           bool
	durations       map[string]time.Duration
	errors          map[string]string
	messages        map[string]string
	PoweredBy       string
	Router          *bear.Mux
	SafeErrorFilter func(error) error
	wares           map[string]bear.HandlerFunc
}

func initDefaults(app *App) {
	app.SetDuration("Cookie", DurationCookie)
	app.SetDuration("Session", DurationSession)
	app.SetError("BadCredentials", ErrorBadCredentials)
	app.SetError("CSRF", ErrorCSRF)
	app.SetError("Generic", ErrorGeneric)
	app.SetError("MethodNotAllowed", ErrorMethodNotAllowed)
	app.SetError("NotFound", ErrorNotFound)
	app.SetError("Parse", ErrorParse)
	app.SetError("Unauthorized", ErrorUnauthorized)
}

func (app *App) Error(key string) string { return app.errors[key] }

func (app *App) SetError(key string, value string) {
	app.errors[key] = value
	InitLog(app, "initialize", fmt.Sprintf("(*forest.App).Error(\"%s\") = %s", key, value))
}

func (app *App) Duration(key string) time.Duration { return app.durations[key] }

func (app *App) SetDuration(key string, value time.Duration) {
	app.durations[key] = value
	InitLog(app, "initialize", fmt.Sprintf("(*forest.App).Duration(\"%s\") = %s", key, value))
}

func (app *App) Message(key string) string { return app.messages[key] }

func (app *App) SetMessage(key string, value string) {
	app.messages[key] = value
	InitLog(app, "initialize", fmt.Sprintf("(*forest.App).Message(\"%s\") = %s", key, value))
}

func (app *App) InstallWare(key string, handler bear.HandlerFunc, message string) error {
	if handler == nil {
		return fmt.Errorf("(*forest.App).InstallWare(\"%s\") was passed a nil handler", key)
	}
	if app.wares[key] != nil {
		message := "overwritten, perhaps multiple Install(Error|Security|Session)Wares invocations"
		println(fmt.Sprintf("(*forest.App).Ware(\"%s\") %s", key, message))
	} else {
		InitLog(app, "install", fmt.Sprintf("(*forest.App).Ware(\"%s\") %s", key, message))
	}
	app.wares[key] = handler
	return nil
}

func (app *App) RegisterRoute(path string, sub SubRouter) { sub.Route(path) }

func (app *App) Response(res http.ResponseWriter, code int, success bool, message string) *Response {
	return &Response{app: app, Code: code, Success: success, Message: message, writer: res}
}

func (app *App) Serve(port string) error {
	if "" == port {
		return fmt.Errorf("forest: no port was specified")
	}
	return http.ListenAndServe(port, app.Router)
}

func (app *App) Ware(key string) bear.HandlerFunc { return app.wares[key] }

func New(debug bool) *App {
	app := &App{
		Debug:     debug,
		durations: make(map[string]time.Duration),
		errors:    make(map[string]string),
		messages:  make(map[string]string),
		Router:    bear.New(),
		wares:     make(map[string]bear.HandlerFunc)}
	initDefaults(app)
	if app.Debug {
		app.Router.Always(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
			ip := req.Header.Get("X-Real-IP")
			if ip == "" {
				ip = req.RemoteAddr
			}
			log.Printf("[%s] %s %s\n", ip, req.Method, req.URL.RequestURI())
			ctx.Next(res, req)
		})
	}
	return app
}
