package wares

import (
	"bytes"
	"encoding/json"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"io/ioutil"
	"net/http"
)

type postBody struct {
	SessionID string `json:"sessionid"`
}

func CSRF(app *forest.App) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		if req.Body == nil {
			app.Response(res, http.StatusBadRequest, forest.Failure, app.Error("CSRF")).Write(nil)
			return
		}
		pb := new(postBody)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Parse"))
			app.Response(res, http.StatusBadRequest, forest.Failure, message).Write(nil)
			return
		}
		// set req.Body back to an untouched io.ReadCloser
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if err := json.Unmarshal(body, pb); err != nil {
			message := app.Error("Parse") + ": " + err.Error()
			app.Response(res, http.StatusBadRequest, forest.Failure, message).Write(nil)
			return
		}
		if sessionID, ok := ctx.Get(forest.SessionID).(string); !ok || sessionID != pb.SessionID {
			app.Response(res, http.StatusBadRequest, forest.Failure, app.Error("CSRF")).Write(nil)
		} else {
			ctx.Next(res, req)
		}
	})
}
