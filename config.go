// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"encoding/json"
	"io/ioutil"
)

const configFile = "forest.json"

type serviceConfig struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type autodiscoveryConfig struct {
	Interface string `json:"interface,omitempty"`
	Name      string `json:"name,omitempty"`
}

type appConfig struct {
	Autodiscovery *autodiscoveryConfig `json:"autodiscovery,omitempty"`
	CookiePath    string
	Debug         bool
	LogRequests   bool
	PoweredBy     string
	ProxyPath     string
	Service       *serviceConfig `json:"service,omitempty"`
}

func loadConfig(app *App) error {
	data, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = json.Unmarshal(data, app.config)
	}
	if app.config.Service == nil {
		app.config.Service = &serviceConfig{}
	}
	if len(app.config.Service.Name) == 0 {
		app.config.Service.Name = "unknown-service"
	}
	if len(app.config.Service.Version) == 0 {
		app.config.Service.Version = "x.x.x"
	}
	if app.config.Autodiscovery == nil {
		app.config.Autodiscovery = &autodiscoveryConfig{}
	}
	return err
}
