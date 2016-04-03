// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

/*
Package forest provides a minimalist framework for writing REST services that
speak JSON.
*/
package forest

import (
	"log"
	"time"
)

const (
	Body                  = "forestbody"
	DurationCookie        = 26 * time.Hour
	DurationSession       = 25 * time.Hour
	ErrorBadCredentials   = "bad credentials"
	ErrorCSRF             = SessionID + " required"
	ErrorGeneric          = "something went wrong"
	ErrorMethodNotAllowed = "method not allowed"
	ErrorNotFound         = "not found"
	ErrorParse            = "json parse error"
	ErrorUnauthorized     = "unauthorized access"
	Failure               = false
	Success               = true
	NoMessage             = ""
	SessionID             = "sessionid"
	SessionRefresh        = "sessionrefresh"
	SessionUser           = "sessionuser"
	SessionUserID         = "sessionuserid"
	Error                 = "foresterror"
	SafeError             = "forestsafeerror"
	UnknownIP             = "unknown-ip"
	UnknownAgent          = "unknown user agent"
	UnknownSession        = "unknown-session"
	WareInstalled         = "forest middleware"
)

func InitLog(app *App, level string, message string) {
	var prefix string
	switch level {
	case "initialize":
		prefix = "[initialized]"
	case "install":
		prefix = "[ installed ]"
	case "listen":
		prefix = "[ listening ]"
	case "warn":
		prefix = "[**warning**]"
	default:
		prefix = "[ undefined ]"
	}
	if app.Debug() {
		log.Printf("%s %s", prefix, message)
	}
}
