package wares

import (
	"fmt"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
	"time"
)

func InstallBodyParser(app *forest.App) {
	installAndLog(app, "BodyParser", BodyParser(app))
}

func InstallErrorWares(app *forest.App) {
	installAndLog(app, "BadRequest", BadRequest(app))
	installAndLog(app, "Conflict", BadRequest(app))
	installAndLog(app, "MethodNotAllowed", MethodNotAllowed(app))
	installAndLog(app, "NotFound", NotFound(app))
	installAndLog(app, "ServerError", ServerError(app))
	installAndLog(app, "Unauthorized", Unauthorized(app))
}

func InstallSecurityWares(app *forest.App) {
	installAndLog(app, "Authenticate", Authenticate(app))
	installAndLog(app, "CSRF", CSRF(app))
}

func InstallSessionWares(app *forest.App, manager SessionManager) {
	installAndLog(app, "SessionDel", SessionDel(app, manager))
	installAndLog(app, "SessionGet", SessionGet(app, manager))
	installAndLog(app, "SessionSet", SessionSet(app, manager))
}

func installAndLog(app *forest.App, key string, value bear.HandlerFunc) {
	message := "forest middleware"
	if err := app.InstallWare(key, value, message); err != nil {
		println(fmt.Sprintf("(*forest.App).Ware(\"%s\") install error: %s", key, err))
	}
}

func safeErrorFilter(app *forest.App, err error, friendly string) error {
	if app.Debug {
		return err
	} else {
		if app.SafeErrorFilter != nil {
			if err := app.SafeErrorFilter(err); err != nil {
				return err
			} else {
				return fmt.Errorf(friendly)
			}
		} else {
			return fmt.Errorf(friendly)
		}
	}
}

func safeErrorMessage(app *forest.App, ctx *bear.Context, friendly string) string {
	if err, ok := ctx.Get(forest.SafeError).(error); ok && err != nil {
		return err.Error()
	} else if err, ok := ctx.Get(forest.Error).(error); ok && err != nil {
		return safeErrorFilter(app, err, friendly).Error()
	} else {
		return friendly
	}
}

func setCookie(res http.ResponseWriter, key string, value string, duration time.Duration) {
	http.SetCookie(res, &http.Cookie{
		Name:     key,
		Value:    value,
		Expires:  time.Now().Add(duration),
		HttpOnly: true,
		MaxAge:   int(duration / time.Second),
		Path:     "/",
		Secure:   true})
}
