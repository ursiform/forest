// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

/*
forest is a micro-framework for building REST services that talk JSON. Its core
unit is a forest.App that is built upon a bear multiplexer for URL routing. It
outputs responses using forest.Response and provides utility methods for many
common tasks required by web services.
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
	WareInstalled         = "forest middleware"
)

func InitLog(app *App, level string, message string) {
	var prefix string
	switch level {
	case "initialize":
		prefix = "[initialized]"
	case "install":
		prefix = "[installed]  "
	default:
		prefix = "[undefined]  "
	}
	if app.Debug {
		log.Printf("%s %s", prefix, message)
	}
}
