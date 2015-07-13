/*
	Package forest
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
