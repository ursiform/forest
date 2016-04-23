// Copyright 2016 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"encoding/json"
	"io/ioutil"
)

const ConfigFile = "bear.json"

type ServiceConfig struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type AppConfig struct {
	CookiePath  string
	Debug       bool
	LogRequests bool
	PoweredBy   string
	ProxyPath   string
	Service     *ServiceConfig `json:"service,omitempty"`
}

func loadConfig(app *App) error {
	data, err := ioutil.ReadFile(ConfigFile)
	if err == nil {
		err = json.Unmarshal(data, app.Config)
	}
	if app.Config.Service == nil {
		app.Config.Service = &ServiceConfig{}
	}
	return err
}
