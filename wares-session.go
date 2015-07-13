// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/ursiform/bear"
	"net/http"
)

func waresSessionDel(app *App, manager SessionManager) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		sessionID, ok := ctx.Get(SessionID).(string)
		if !ok {
			ctx.Set(Error, fmt.Errorf("SessionDel %s: %v", SessionID, ctx.Get(SessionID)))
			message := app.safeErrorMessage(ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, Failure, message).Write(nil)
			return
		}
		userID, ok := ctx.Get(SessionUserID).(string)
		if !ok {
			ctx.Set(Error, fmt.Errorf("SessionDel %s: %v",
				SessionUserID, ctx.Get(SessionUserID)))
			message := app.safeErrorMessage(ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, Failure, message).Write(nil)
			return
		}
		if err := manager.Delete(sessionID, userID); err != nil {
			ctx.Set(Error, err)
			message := app.safeErrorMessage(ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, Failure, message).Write(nil)
			return
		}
		ctx.Next(res, req)
	})
}

func waresSessionGet(app *App, manager SessionManager) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		createEmptySession := func(sessionID string) {
			path := app.CookiePath
			if len(path) == 0 {
				path = "/"
			}
			key := SessionID
			value := sessionID
			duration := app.Duration("Cookie")
			app.SetCookie(res, path, key, value, duration) // reset the cookie
			manager.CreateEmpty(sessionID, ctx)
			ctx.Next(res, req)
		}
		cookie, err := req.Cookie(SessionID)
		if err != nil || cookie.Value == "" {
			createEmptySession(uuid.New())
			return
		}
		sessionID := cookie.Value
		userID, userJSON, err := manager.Read(sessionID)
		if err != nil || userID == "" || userJSON == "" {
			createEmptySession(sessionID)
			return
		}
		if err := manager.Create(sessionID, userID, userJSON, ctx); err != nil {
			println(fmt.Sprintf("error creating session: %s", err))
			defer func(sessionID string, userID string) {
				if err := manager.Delete(sessionID, userID); err != nil {
					println(fmt.Sprintf("error deleting session: %s", err))
				}
			}(sessionID, userID)
			createEmptySession(sessionID)
			return
		}
		// if SessionRefresh is set to false, the session will not refresh,
		// otherwise, if it is not set or if it is set to true, the session is refreshed
		if refresh, ok := ctx.Get(SessionRefresh).(bool); !ok || refresh {
			path := app.CookiePath
			if len(path) == 0 {
				path = "/"
			}
			key := SessionID
			value := sessionID
			duration := app.Duration("Cookie")
			app.SetCookie(res, path, key, value, duration) // refreshes the cookie
			defer func(sessionID string, userJSON string) {
				if err := manager.Update(sessionID, userID, userJSON, app.Duration("Session")); err != nil {
					println(fmt.Sprintf("error updating session: %s", err))
				}
			}(sessionID, userJSON)
		}
		ctx.Next(res, req)
	})
}

func waresSessionSet(app *App, manager SessionManager) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		userJSON, err := manager.Marshal(ctx)
		if err != nil {
			ctx.Set(Error, err)
			message := app.safeErrorMessage(ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, Failure, message).Write(nil)
			return
		}
		sessionID, ok := ctx.Get(SessionID).(string)
		if !ok {
			ctx.Set(Error, fmt.Errorf("%s: %v", SessionID, ctx.Get(SessionID)))
			message := app.safeErrorMessage(ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, Failure, message).Write(nil)
			return
		}
		userID, ok := ctx.Get(SessionUserID).(string)
		if !ok {
			ctx.Set(Error, fmt.Errorf("%s: %v", SessionUserID, ctx.Get(SessionUserID)))
			message := app.safeErrorMessage(ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, Failure, message).Write(nil)
			return
		}
		if err := manager.Update(sessionID, userID, string(userJSON), app.Duration("Session")); err != nil {
			ctx.Set(Error, err)
			message := app.safeErrorMessage(ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, Failure, message).Write(nil)
			return
		}
		ctx.Next(res, req)
	})
}
