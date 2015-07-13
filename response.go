package forest

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	app     *App
	Code    int         `json:"-"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	writer  http.ResponseWriter
}

func (res *Response) Write(data interface{}) (bytes int, err error) {
	res.Data = data
	output, _ := json.Marshal(res)
	res.writer.Header().Set("Content-Type", "application/json")
	if 0 < len(res.app.PoweredBy) {
		res.writer.Header().Set("X-Powered-By", res.app.PoweredBy)
	}
	res.writer.WriteHeader(res.Code)
	return res.writer.Write(output)
}
