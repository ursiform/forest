// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
)

func ErrorsBadRequest(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(res, http.StatusBadRequest, forest.Failure, message).Write(nil)
	})
}

func ErrorsConflict(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(res, http.StatusConflict, forest.Failure, message).Write(nil)
	})
}

func ErrorsMethodNotAllowed(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("MethodNotAllowed"))
		app.Response(res, http.StatusMethodNotAllowed, forest.Failure, message).Write(nil)
	})
}

func ErrorsNotFound(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("NotFound"))
		app.Response(res, http.StatusNotFound, forest.Failure, message).Write(nil)
	})
}

func ErrorsServerError(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Generic"))
		app.Response(res, http.StatusInternalServerError, forest.Failure, message).Write(nil)
	})
}

func ErrorsUnauthorized(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		message := safeErrorMessage(app, ctx, app.Error("Unauthorized"))
		app.Response(res, http.StatusUnauthorized, forest.Failure, message).Write(nil)
	})
}
