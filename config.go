// Copyright 2016 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ursiform/logger"
)

const ConfigFile = "bear.json"

type ServiceConfig struct {
	Address      string `json:"address,omitempty"`
	LogLevelName string `json:"loglevel,omitempty"`
	LogRequests  bool   `json:"logrequests,omitempty"`
	Name         string `json:"name,omitempty"`
	Version      string `json:"version,omitempty"`
}

type AppConfig struct {
	CookiePath string
	LogLevel   int
	PoweredBy  string
	Debug      bool           `json:"debug,omitempty"`
	Service    *ServiceConfig `json:"service,omitempty"`
}

func loadConfig(app *App) error {
	data, err := ioutil.ReadFile(ConfigFile)
	if err == nil {
		err = json.Unmarshal(data, app.Config)
	}
	if app.Config.Service == nil {
		app.Config.Service = &ServiceConfig{}
	}
	if len(app.Config.Service.LogLevelName) == 0 {
		app.Config.Service.LogLevelName = "listen"
	}
	level, ok := logger.LogLevel[app.Config.Service.LogLevelName]
	if !ok {
		message := fmt.Sprintf("%s=\"%s\" in %s is invalid; using \"%s\"",
			"service.loglevel", app.Config.Service.LogLevelName, ConfigFile, "debug")
		app.Config.Service.LogLevelName = "debug"
		app.Config.LogLevel = logger.Debug
		logger.MustLog(logger.Error, message)
	} else {
		app.Config.LogLevel = level
	}
	return err
}
