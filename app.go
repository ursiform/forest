// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ursiform/bear"
	"github.com/ursiform/logger"
)

const address = ":80"

type App struct {
	*bear.Mux
	Config          *AppConfig
	durations       map[string]time.Duration
	errors          map[string]string
	logger          *logger.Logger
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

// Duration gets the duration for a specific key, e.g. "Cookie" expiration.
func (app *App) Duration(key string) time.Duration { return app.durations[key] }

// SetDuration sets the duration for a specific key, e.g. "Cookie" expiration.
func (app *App) SetDuration(key string, value time.Duration) {
	app.durations[key] = value
	app.Log(logger.Init, fmt.Sprintf("Duration(\"%s\") = %s", key, value))
}

// Error gets the error for a specific key, e.g. "Unauthorized".
func (app *App) Error(key string) string { return app.errors[key] }

// SetError sets the error for a specific key, e.g. "Unauthorized".
func (app *App) SetError(key string, value string) {
	app.errors[key] = value
	app.Log(logger.Init, fmt.Sprintf("Error(\"%s\") = %s", key, value))
}

// Log outputs a log message at a specified log level.
func (app *App) Log(level int, message string) {
	app.logger.Log(level, message)
}

// Message gets the app message for a specific key, e.g. "AlreadyLoggedIn".
func (app *App) Message(key string) string { return app.messages[key] }

// SetMessage sets the app message for a specific key, e.g. "AlreadyLoggedIn".
func (app *App) SetMessage(key string, value string) {
	app.messages[key] = value
	app.Log(logger.Init, fmt.Sprintf("Message(\"%s\") = %s", key, value))
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

func (app *App) InstallWare(
	key string, handler func(ctx *bear.Context), message string) error {
	if handler == nil {
		return fmt.Errorf("InstallWare(\"%s\") is nil", key)
	}
	if app.wares[key] != nil {
		app.Log(logger.Warn, fmt.Sprintf("Ware(\"%s\") %s",
			key, "overwritten, perhaps multiple Install* invocations"))
	} else {
		app.Log(logger.Install, fmt.Sprintf("Ware(\"%s\") %s", key, message))
	}
	app.wares[key] = handler
	return nil
}

func (app *App) ListenAndServe() error {
	return http.ListenAndServe(app.Config.Service.Address, app.Mux)
}

func (app *App) ListenAndServeTLS(certFile, keyFile string) error {
	addr := app.Config.Service.Address
	return http.ListenAndServeTLS(addr, certFile, keyFile, app.Mux)
}

func (app *App) On(verb string, pattern string, handlers ...interface{}) error {
	app.Log(logger.Init, fmt.Sprintf("%s %s", verb, pattern))
	return app.Mux.On(verb, pattern, handlers...)
}

func (app *App) RegisterRoute(path string, sub SubRouter) { sub.Route(path) }

func (app *App) Response(ctx *bear.Context,
	code int, success bool, message string) *Response {
	return &Response{
		app:     app,
		ctx:     ctx,
		code:    code,
		Success: success,
		Message: message}
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

func New() *App {
	app := new(App)
	app.Config = new(AppConfig)
	app.Mux = bear.New()
	err := loadConfig(app)
	app.logger, _ = logger.New(app.Config.LogLevel)
	if err != nil {
		logger.MustLog(logger.Warn, ConfigFile+" was not loaded")
	}
	if app.Config.Service.Address == "" {
		message := fmt.Sprintf("%s is not defined in %s, using default %s",
			"service.address", ConfigFile, address)
		app.Config.Service.Address = address
		logger.MustLog(logger.Warn, message)
	}
	app.durations = make(map[string]time.Duration)
	app.errors = make(map[string]string)
	app.messages = make(map[string]string)
	app.wares = make(map[string]func(ctx *bear.Context))
	initDefaults(app)
	return app
}
