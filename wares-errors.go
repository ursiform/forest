package forest

import (
	"github.com/ursiform/bear"
	"net/http"
)

func waresErrorsBadRequest(app *App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := app.safeErrorMessage(ctx, app.Error("Generic"))
		app.Response(res, http.StatusBadRequest, Failure, message).Write(nil)
	})
}

func waresErrorsConflict(app *App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := app.safeErrorMessage(ctx, app.Error("Generic"))
		app.Response(res, http.StatusConflict, Failure, message).Write(nil)
	})
}

func waresErrorsMethodNotAllowed(app *App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := app.safeErrorMessage(ctx, app.Error("MethodNotAllowed"))
		app.Response(res, http.StatusMethodNotAllowed, Failure, message).Write(nil)
	})
}

func waresErrorsNotFound(app *App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := app.safeErrorMessage(ctx, app.Error("NotFound"))
		app.Response(res, http.StatusNotFound, Failure, message).Write(nil)
	})
}

func waresErrorsServerError(app *App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := app.safeErrorMessage(ctx, app.Error("Generic"))
		app.Response(res, http.StatusInternalServerError, Failure, message).Write(nil)
	})
}

func waresErrorsUnauthorized(app *App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := app.safeErrorMessage(ctx, app.Error("Unauthorized"))
		app.Response(res, http.StatusUnauthorized, Failure, message).Write(nil)
	})
}
