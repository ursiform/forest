// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

/*
Package wares is a collection of bear.HandlerFunc middleware generators for use
with a forest.App instance.
*/
package wares

import (
	"fmt"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
)

func InstallBodyParser(app *forest.App) {
	app.InstallWare("BodyParser", BodyParser(app), forest.WareInstalled)
}

func InstallErrorWares(app *forest.App) {
	app.InstallWare("BadRequest", ErrorsBadRequest(app), forest.WareInstalled)
	app.InstallWare("Conflict", ErrorsBadRequest(app), forest.WareInstalled)
	app.InstallWare("MethodNotAllowed", ErrorsMethodNotAllowed(app), forest.WareInstalled)
	app.InstallWare("NotFound", ErrorsNotFound(app), forest.WareInstalled)
	app.InstallWare("ServerError", ErrorsServerError(app), forest.WareInstalled)
	app.InstallWare("Unauthorized", ErrorsUnauthorized(app), forest.WareInstalled)
}

func InstallSecurityWares(app *forest.App) {
	app.InstallWare("Authenticate", Authenticate(app), forest.WareInstalled)
	app.InstallWare("CSRF", CSRF(app), forest.WareInstalled)
}

func InstallSessionWares(app *forest.App, manager SessionManager) {
	app.InstallWare("SessionDel", SessionDel(app, manager), forest.WareInstalled)
	app.InstallWare("SessionGet", SessionGet(app, manager), forest.WareInstalled)
	app.InstallWare("SessionSet", SessionSet(app, manager), forest.WareInstalled)
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
