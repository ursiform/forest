package wares

import (
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
)

func BadRequest(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(res, http.StatusBadRequest, forest.Failure, message).Write(nil)
	})
}

func Conflict(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(res, http.StatusConflict, forest.Failure, message).Write(nil)
	})
}

func MethodNotAllowed(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("MethodNotAllowed"))
		app.Response(res, http.StatusMethodNotAllowed, forest.Failure, message).Write(nil)
	})
}

func NotFound(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("NotFound"))
		app.Response(res, http.StatusNotFound, forest.Failure, message).Write(nil)
	})
}

func ServerError(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(res, http.StatusInternalServerError, forest.Failure, message).Write(nil)
	})
}

func Unauthorized(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Unauthorized"))
		app.Response(res, http.StatusUnauthorized, forest.Failure, message).Write(nil)
	})
}
