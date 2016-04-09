// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"fmt"
	"github.com/ursiform/bear"
	"net/http"
	"time"
)

type App struct {
	*bear.Mux
	Config          *AppConfig
	durations       map[string]time.Duration
	errors          map[string]string
	messages        map[string]string
	SafeErrorFilter func(error) error
	wares           map[string]func(ctx *bear.Context)
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

// Getters and setters

// CookiePath gets the cookie path for cookies the app sets.
func (app *App) CookiePath() string {
	if len(app.Config.CookiePath) > 0 {
		return app.Config.CookiePath
	} else {
		return ""
	}
}

// SetCookiePath sets the cookie path for cookies the app sets.
func (app *App) SetCookiePath(value string) {
	app.Config.CookiePath = value
}

// Debug gets the app debug flag.
func (app *App) Debug() bool { return app.Config.Debug }

// SetDebug sets the app debug flag.
func (app *App) SetDebug(value bool) { app.Config.Debug = value }

// Duration gets the duration for a specific key, e.g. "Cookie" expiration.
func (app *App) Duration(key string) time.Duration { return app.durations[key] }

// SetDuration sets the duration for a specific key, e.g. "Cookie" expiration.
func (app *App) SetDuration(key string, value time.Duration) {
	app.durations[key] = value
	output := fmt.Sprintf("Duration(\"%s\") = %s", key, value)
	InitLog(app, "initialize", output)
}

// Error gets the error for a specific key, e.g. "Unauthorized".
func (app *App) Error(key string) string { return app.errors[key] }

// SetError sets the error for a specific key, e.g. "Unauthorized".
func (app *App) SetError(key string, value string) {
	app.errors[key] = value
	output := fmt.Sprintf("Error(\"%s\") = %s", key, value)
	InitLog(app, "initialize", output)
}

// LogRequests gets the app request logging flag.
func (app *App) LogRequests() bool { return app.Config.LogRequests }

// SetLogRequests sets the app request logging flag.
func (app *App) SetLogRequests(value bool) { app.Config.LogRequests = value }

// Message gets the app message for a specific key, e.g. "AlreadyLoggedIn".
func (app *App) Message(key string) string { return app.messages[key] }

// SetMessage sets the app message for a specific key, e.g. "AlreadyLoggedIn".
func (app *App) SetMessage(key string, value string) {
	app.messages[key] = value
	output := fmt.Sprintf("Message(\"%s\") = %s", key, value)
	InitLog(app, "initialize", output)
}

// PoweredBy gets the response X-Powered-By HTTP header.
func (app *App) PoweredBy() string {
	if len(app.Config.PoweredBy) > 0 {
		return app.Config.PoweredBy
	} else {
		return ""
	}
}

// SetPoweredBy sets the response X-Powered-By HTTP header.
func (app *App) SetPoweredBy(value string) {
	app.Config.PoweredBy = value
}

// ProxyPath gets the reverse proxy path for self-documentation.
func (app *App) ProxyPath() string {
	if len(app.Config.ProxyPath) > 0 {
		return app.Config.ProxyPath
	} else {
		return ""
	}
}

// SetProxyPath sets the reverse proxy path for self-documentation.
func (app *App) SetProxyPath(proxyPath string) {
	app.Config.ProxyPath = proxyPath
}

func (app *App) InstallWare(
	key string, handler func(ctx *bear.Context), message string) error {
	if handler == nil {
		return fmt.Errorf("InstallWare(\"%s\") is nil", key)
	}
	if app.wares[key] != nil {
		output := "overwritten, perhaps multiple Install* invocations"
		println(fmt.Sprintf("Ware(\"%s\") %s", key, output))
	} else {
		output := fmt.Sprintf("Ware(\"%s\") %s", key, message)
		InitLog(app, "install", output)
	}
	app.wares[key] = handler
	return nil
}

func (app *App) On(verb string, pattern string, handlers ...interface{}) error {
	InitLog(app, "listen", fmt.Sprintf("%s %s%s", verb, app.ProxyPath(), pattern))
	return app.Mux.On(verb, pattern, handlers...)
}

func (app *App) RegisterRoute(path string, sub SubRouter) { sub.Route(path) }

func (app *App) Response(ctx *bear.Context,
	code int, success bool, message string) *Response {
	return &Response{
		app:     app,
		ctx:     ctx,
		Code:    code,
		Success: success,
		Message: message}
}

func (app *App) Serve(port string) error {
	if "" == port {
		return fmt.Errorf("forest: no port was specified")
	}
	return http.ListenAndServe(port, app.Mux)
}

func (app *App) SetCookie(
	ctx *bear.Context, path, key, value string, duration time.Duration) {
	response := &Response{app: app, ctx: ctx}
	response.SetCookie(path, key, value, duration)
}

func (app *App) Ware(key string) func(ctx *bear.Context) {
	if handler := app.wares[key]; handler != nil {
		return handler
	}
	return func(ctx *bear.Context) {
		app.Response(
			ctx,
			http.StatusInternalServerError,
			Failure,
			fmt.Sprintf("Ware(%s) is nil", key)).Write(nil)
	}
}

func New(debug bool) *App {
	app := new(App)
	app.Config = new(AppConfig)
	app.SetDebug(debug) // Set debug before using InitLog.
	app.Mux = bear.New()
	if err := loadConfig(app); err != nil {
		InitLog(app, "warning", ConfigFile+" was not loaded")
	}
	app.durations = make(map[string]time.Duration)
	app.errors = make(map[string]string)
	app.messages = make(map[string]string)
	app.wares = make(map[string]func(ctx *bear.Context))
	initDefaults(app)
	return app
}
