// Copyright 2016 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ursiform/logger"
)

const ConfigFile = "forest.json"

type Config struct {
	Address      string `json:"address,omitempty"`
	CookiePath   string
	File         string
	LogLevel     int
	LogLevelName string `json:"loglevel,omitempty"`
	LogRequests  bool   `json:"logrequests,omitempty"`
	Name         string `json:"name,omitempty"`
	PoweredBy    string
	Debug        bool   `json:"debug,omitempty"`
	Version      string `json:"version,omitempty"`
}

func loadConfig(app *App) error {
	data, err := ioutil.ReadFile(app.Config.File)
	if err == nil {
		err = json.Unmarshal(data, app.Config)
	}
	if len(app.Config.LogLevelName) == 0 {
		app.Config.LogLevelName = "listen"
	}
	level, ok := logger.LogLevel[app.Config.LogLevelName]
	if !ok {
		logger.MustError("loglevel=\"%s\" in %s is invalid; using \"%s\"",
			app.Config.LogLevelName, app.Config.File, "debug")
		app.Config.LogLevelName = "debug"
		app.Config.LogLevel = logger.Debug
	} else {
		app.Config.LogLevel = level
	}
	return err
}
