// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"encoding/json"
	"github.com/ursiform/bear"
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
	return res.ctx.ResponseWriter.Write(output)
}
