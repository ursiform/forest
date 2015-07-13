package forest

import (
	"bytes"
	"encoding/json"
	"github.com/ursiform/bear"
	"io/ioutil"
	"net/http"
)

type csrfPostBody struct {
	SessionID string `json:"sessionid"` // forest.SessionID == "sessionid"
}

func waresCSRF(app *App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		if req.Body == nil {
			app.Response(res, http.StatusBadRequest, Failure, app.Error("CSRF")).Write(nil)
			return
		}
		pb := new(csrfPostBody)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			ctx.Set(Error, err)
			message := app.safeErrorMessage(ctx, app.Error("Parse"))
			app.Response(res, http.StatusBadRequest, Failure, message).Write(nil)
			return
		}
		// set req.Body back to an untouched io.ReadCloser
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if err := json.Unmarshal(body, pb); err != nil {
			message := app.Error("Parse") + ": " + err.Error()
			app.Response(res, http.StatusBadRequest, Failure, message).Write(nil)
			return
		}
		if sessionID, ok := ctx.Get(SessionID).(string); !ok || sessionID != pb.SessionID {
			app.Response(res, http.StatusBadRequest, Failure, app.Error("CSRF")).Write(nil)
		} else {
			ctx.Next(res, req)
		}
	})
}
