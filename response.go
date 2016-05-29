// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ursiform/bear"
	"github.com/ursiform/logger"
)

type Response struct {
	app     *App
	code    int
	ctx     *bear.Context
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
}

func (res *Response) SetCookie(
	path, key, value string, duration time.Duration) {
	http.SetCookie(res.ctx.ResponseWriter, &http.Cookie{
		Name:     key,
		Value:    value,
		Expires:  time.Now().Add(duration),
		HttpOnly: true,
		MaxAge:   int(duration / time.Second),
		Path:     path,
		Secure:   true})
}

func (res *Response) Write(data interface{}) (bytes int, err error) {
	res.Data = data
	output, _ := json.Marshal(res)
	res.ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	poweredBy := res.app.Config.PoweredBy
	if 0 < len(poweredBy) {
		res.ctx.ResponseWriter.Header().Set("X-Powered-By", poweredBy)
	}
	res.ctx.ResponseWriter.WriteHeader(res.code)
	bytes, err = res.ctx.ResponseWriter.Write(output)
	if !res.app.Config.LogRequests &&
		!res.app.Config.Debug &&
		res.app.Config.LogLevel < logger.Request {
		return
	}
	// First, try the X-Real-IP header from the reverse proxy.
	ip := res.ctx.Request.Header.Get("X-Real-IP")
	// If X-Real-IP does not exist, try the REMOTE_ADDR.
	if ip == "" {
		ip = res.ctx.Request.RemoteAddr
	}
	// If IP is still blank (perhaps in tests), display UnknownIP.
	if ip == "" {
		ip = UnknownIP
	}
	method := res.ctx.Request.Method
	uri := res.ctx.Request.URL.RequestURI()
	proto := res.ctx.Request.Proto
	code := res.code
	agent := res.ctx.Request.UserAgent()
	if agent == "" {
		agent = UnknownAgent
	}
	sessionID, ok := res.ctx.Get(SessionID).(string)
	message := res.Message
	if message == NoMessage {
		message = "no message"
	}
	if !ok || sessionID == "" {
		sessionID = UnknownSession
	}
	logger.MustRequest("%s@%s \"%s %s %s\" %d %d \"%s\" [%s]",
		sessionID, ip, method, uri, proto, code, bytes, agent, message)
	return
}
