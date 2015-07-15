// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest_test

import (
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
)

type responseFormat struct {
	Foo string `json:"foo"`
}

type router struct{ *forest.App }

func (app *router) respondSuccess(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
	data := &responseFormat{Foo: "foo"}
	app.Response(res, http.StatusOK, forest.Success, forest.NoMessage).Write(data)
}

func (app *router) setCookie(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
	path := "/"
	cookieName := "foo"
	cookieValue := "Foo"
	app.SetCookie(res, path, cookieName, cookieValue, app.Duration("Cookie"))
	ctx.Next(res, req)
}

func (app *router) Route(path string) {
	app.Router.On("GET", path, app.setCookie, app.respondSuccess)
	app.Router.On("*", path, app.Ware("MethodNotAllowed"))
}

func newRouter(parent *forest.App) *router { return &router{parent} }
