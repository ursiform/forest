// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"fmt"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
)

func BodyParser(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		destination, ok := ctx.Get(forest.Body).(forest.Populater)
		if !ok {
			ctx.Set(forest.Error, fmt.Errorf("(*forest.App).BodyParser"))
			message := safeErrorMessage(app, ctx, app.Error("Parse"))
			app.Response(res, http.StatusBadRequest, forest.Failure, message).Write(nil)
			return
		}
		if req.Body == nil {
			ctx.Set(forest.SafeError, fmt.Errorf("%s: body is empty", app.Error("Parse")))
			message := safeErrorMessage(app, ctx, app.Error("Parse"))
			app.Response(res, http.StatusBadRequest, forest.Failure, message).Write(nil)
			return
		}
		if err := destination.Populate(req.Body); err != nil {
			ctx.Set(forest.SafeError, fmt.Errorf("%s: %s", app.Error("Parse"), err))
			message := safeErrorMessage(app, ctx, app.Error("Parse"))
			app.Response(res, http.StatusBadRequest, forest.Failure, message).Write(nil)
			return
		}
		ctx.Next(res, req)
	})
}
