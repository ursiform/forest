// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package wares

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"net/http"
)

func SessionDel(app *forest.App, manager SessionManager) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		sessionID, ok := ctx.Get(forest.SessionID).(string)
		if !ok {
			ctx.Set(forest.Error, fmt.Errorf("SessionDel %s: %v", forest.SessionID, ctx.Get(forest.SessionID)))
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, forest.Failure, message).Write(nil)
			return
		}
		userID, ok := ctx.Get(forest.SessionUserID).(string)
		if !ok {
			ctx.Set(forest.Error, fmt.Errorf("SessionDel %s: %v",
				forest.SessionUserID, ctx.Get(forest.SessionUserID)))
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, forest.Failure, message).Write(nil)
			return
		}
		if err := manager.Delete(sessionID, userID); err != nil {
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, forest.Failure, message).Write(nil)
			return
		}
		ctx.Next(res, req)
	})
}

func SessionGet(app *forest.App, manager SessionManager) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		cookieName := forest.SessionID
		createEmptySession := func(sessionID string) {
			cookiePath := app.CookiePath
			if cookiePath == "" {
				cookiePath = "/"
			}
			cookieValue := sessionID
			cookieDuration := app.Duration("Cookie")
			app.SetCookie(res, cookiePath, cookieName, cookieValue, cookieDuration) // reset the cookie
			manager.CreateEmpty(sessionID, ctx)
			ctx.Next(res, req)
		}
		cookie, err := req.Cookie(cookieName)
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
		if refresh, ok := ctx.Get(forest.SessionRefresh).(bool); !ok || refresh {
			cookiePath := app.CookiePath
			if cookiePath == "" {
				cookiePath = "/"
			}
			cookieName := forest.SessionID
			cookieValue := sessionID
			cookieDuration := app.Duration("Cookie")
			app.SetCookie(res, cookiePath, cookieName, cookieValue, cookieDuration) // refreshes the cookie
			defer func(sessionID string, userJSON string) {
				if err := manager.Update(sessionID, userID, userJSON, app.Duration("Session")); err != nil {
					println(fmt.Sprintf("error updating session: %s", err))
				}
			}(sessionID, userJSON)
		}
		ctx.Next(res, req)
	})
}

func SessionSet(app *forest.App, manager SessionManager) bear.HandlerFunc {
	return bear.HandlerFunc(func(res http.ResponseWriter, req *http.Request, ctx *bear.Context) {
		userJSON, err := manager.Marshal(ctx)
		if err != nil {
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, forest.Failure, message).Write(nil)
			return
		}
		sessionID, ok := ctx.Get(forest.SessionID).(string)
		if !ok {
			ctx.Set(forest.Error, fmt.Errorf("%s: %v", forest.SessionID, ctx.Get(forest.SessionID)))
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, forest.Failure, message).Write(nil)
			return
		}
		userID, ok := ctx.Get(forest.SessionUserID).(string)
		if !ok {
			ctx.Set(forest.Error, fmt.Errorf("%s: %v", forest.SessionUserID, ctx.Get(forest.SessionUserID)))
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, forest.Failure, message).Write(nil)
			return
		}
		if err := manager.Update(sessionID, userID, string(userJSON), app.Duration("Session")); err != nil {
			ctx.Set(forest.Error, err)
			message := safeErrorMessage(app, ctx, app.Error("Generic"))
			app.Response(res, http.StatusInternalServerError, forest.Failure, message).Write(nil)
			return
		}
		ctx.Next(res, req)
	})
}
