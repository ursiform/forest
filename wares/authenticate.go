package wares

import (
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
)

func Authenticate(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		userID, ok := ctx.Get(forest.SessionUserID).(string)
		if !ok || len(userID) == 0 {
			app.Response(res, http.StatusUnauthorized, forest.Failure, app.Error("Unauthorized")).Write(nil)
			return
		}
		ctx.Next(res, req)
	})
}
