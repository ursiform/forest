package forest

import (
	"fmt"
	"github.com/ursiform/bear"
	"net/http"
)

func waresBodyParser(app *App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		destination, ok := ctx.Get(Body).(Populater)
		if !ok {
			ctx.Set(Error, fmt.Errorf("(*forest.App).ParseBody"))
			message := app.safeErrorMessage(ctx, app.Error("Parse"))
			app.Response(res, http.StatusBadRequest, Failure, message).Write(nil)
			return
		}
		if req.Body == nil {
			ctx.Set(SafeError, fmt.Errorf("%s: body is empty", app.Error("Parse")))
			message := app.safeErrorMessage(ctx, app.Error("Parse"))
			app.Response(res, http.StatusBadRequest, Failure, message).Write(nil)
			return
		}
		if err := destination.Populate(req.Body); err != nil {
			ctx.Set(SafeError, fmt.Errorf("%s: %s", app.Error("Parse"), err))
			message := app.safeErrorMessage(ctx, app.Error("Parse"))
			app.Response(res, http.StatusBadRequest, Failure, message).Write(nil)
			return
		}
		ctx.Next(res, req)
	})
}
