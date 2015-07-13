package forest

import (
	"github.com/ursiform/bear"
	"net/http"
)

func waresAuthenticate(app *App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		userID, ok := ctx.Get(SessionUserID).(string)
		if !ok || len(userID) == 0 {
			app.Response(res, http.StatusUnauthorized, Failure, app.Error("Unauthorized")).Write(nil)
			return
		}
		ctx.Next(res, req)
	})
}
