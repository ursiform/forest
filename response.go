// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"encoding/json"
	"github.com/ursiform/bear"
	"log"
	"net/http"
	"time"
)

type Response struct {
	app     *App
	Code    int `json:"-"`
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
	if 0 < len(res.app.PoweredBy) {
		res.ctx.ResponseWriter.Header().Set("X-Powered-By", res.app.PoweredBy)
	}
	res.ctx.ResponseWriter.WriteHeader(res.Code)
	bytes, err = res.ctx.ResponseWriter.Write(output)
	if !res.app.Debug {
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
	code := res.Code
	agent := res.ctx.Request.UserAgent()
	if agent == "" {
		agent = UnknownAgent
	}
	sessionID, ok := res.ctx.Get(SessionID).(string)
	message := res.Message
	if !ok || sessionID == "" {
		sessionID = UnknownSession
	}
	log.Printf("[%s/%s] \"%s %s %s\" %d %d \"%s\" %s\n",
		ip, sessionID, method, uri, proto, code, bytes, agent, message)
	return
}
