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
	CookiePath      string
	Debug           bool
	durations       map[string]time.Duration
	errors          map[string]string
	LogRequests     bool
	messages        map[string]string
	PoweredBy       string
	ProxyPath       string
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

// Getters and setters

// Duration gets the duration for a specific key, e.g. "Cookie" expiration.
func (app *App) Duration(key string) time.Duration { return app.durations[key] }

// SetDuration sets the duration for a specific key, e.g. "Cookie" expiration.
func (app *App) SetDuration(key string, value time.Duration) {
	app.durations[key] = value
	output := fmt.Sprintf("(*forest.App).Duration(\"%s\") = %s", key, value)
	InitLog(app, "initialize", output)
}

// Error gets the error for a specific key, e.g. "Unauthorized".
func (app *App) Error(key string) string { return app.errors[key] }

// SetError sets the error for a specific key, e.g. "Unauthorized".
func (app *App) SetError(key string, value string) {
	app.errors[key] = value
	output := fmt.Sprintf("(*forest.App).Error(\"%s\") = %s", key, value)
	InitLog(app, "initialize", output)
}

// Message gets the app message for a specific key, e.g. "AlreadyLoggedIn".
func (app *App) Message(key string) string { return app.messages[key] }

// SetMessage sets the app message for a specific key, e.g. "AlreadyLoggedIn".
func (app *App) SetMessage(key string, value string) {
	app.messages[key] = value
	output := fmt.Sprintf("(*forest.App).Message(\"%s\") = %s", key, value)
	InitLog(app, "initialize", output)
}

func (app *App) InstallWare(key string,
	handler bear.HandlerFunc, message string) error {
	if handler == nil {
		return fmt.Errorf("(*forest.App).InstallWare(\"%s\") is nil", key)
	}
	if app.wares[key] != nil {
		output := "overwritten, perhaps multiple Install* invocations"
		println(fmt.Sprintf("(*forest.App).Ware(\"%s\") %s", key, output))
	} else {
		output := fmt.Sprintf("(*forest.App).Ware(\"%s\") %s", key, message)
		InitLog(app, "install", output)
	}
	app.wares[key] = handler
	return nil
}

func (app *App) On(verb string, pattern string, handlers ...interface{}) error {
	InitLog(app, "listen", fmt.Sprintf("%s %s%s", verb, app.ProxyPath, pattern))
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

func (app *App) Ware(key string) bear.HandlerFunc {
	handler := app.wares[key]
	if handler != nil {
		return handler
	}
	errorHandler := func(
		_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
		message := fmt.Sprintf("(*forest.App).Ware(%s) is nil", key)
		app.Response(ctx,
			http.StatusInternalServerError, Failure, message).Write(nil)
	}
	return bear.HandlerFunc(errorHandler)
}

func New(debug bool) *App {
	app := new(App)
	app.Mux = bear.New()
	app.Debug = debug
	app.durations = make(map[string]time.Duration)
	app.errors = make(map[string]string)
	app.messages = make(map[string]string)
	app.wares = make(map[string]bear.HandlerFunc)
	initDefaults(app)
	return app
}
