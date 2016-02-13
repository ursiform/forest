// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest_test

import (
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"github.com/ursiform/forest-wares"
	"net/http"
)

type responseFormat struct {
	Foo string `json:"foo"`
}

type router struct{ *forest.App }

func (app *router) respondSuccess(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	data := &responseFormat{Foo: "foo"}
	app.Response(ctx,
		http.StatusOK, forest.Success, forest.NoMessage).Write(data)
}

func (app *router) setCookie(
	_ http.ResponseWriter, _ *http.Request, ctx *bear.Context) {
	path := "/"
	cookieName := "foo"
	cookieValue := "bar"
	app.SetCookie(ctx, path, cookieName, cookieValue, app.Duration("Cookie"))
	ctx.Next()
}

func (app *router) Route(path string) {
	app.On("GET", path+"/nonexistent",
		app.Ware("NonExistent"), app.respondSuccess)
	app.On("GET", path,
		app.setCookie, app.respondSuccess)
	app.On("*", path,
		app.Ware("MethodNotAllowed"))
}

func newRouter(parent *forest.App) *router {
	wares.InstallErrorWares(parent)
	return &router{parent}
}
