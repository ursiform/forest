// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares_test

import (
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"github.com/ursiform/forest/wares"
	"net/http"
)

type responseFormat struct {
	Foo string `json:"foo"`
}

type router struct{ *forest.App }

func (app *router) authenticate(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
	ctx.Set(forest.SessionID, "some session id").Set(forest.SessionUserID, "some user id").Next(res, req)
}

func (app *router) respondSuccess(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
	data := &responseFormat{Foo: "foo"}
	println(req.URL.Path)
	app.Response(res, http.StatusOK, forest.Success, forest.NoMessage).Write(data)
}

func (app *router) Route(path string) {
	app.Router.On("GET", path, app.respondSuccess)
	app.Router.On("GET", path+"/authenticate/failure", app.Ware("Authenticate"), app.respondSuccess)
	app.Router.On("GET", path+"/authenticate/success", app.authenticate, app.Ware("Authenticate"), app.respondSuccess)
	app.Router.On("DELETE", path, app.Ware("Unauthorized"))
	app.Router.On("*", path, app.Ware("MethodNotAllowed"))
}

func newRouter(parent *forest.App) *router {
	wares.InstallErrorWares(parent)
	wares.InstallSecurityWares(parent)
	return &router{parent}
}
