// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"encoding/json"
	"io/ioutil"
)

const configFile = "forest.json"

type disqueConfig struct {
	Hosts []string `json:"hosts,omitempty"`
}

type appConfig struct {
	CookiePath  string
	Debug       bool
	Disque      *disqueConfig `json:"disque,omitempty"`
	LogRequests bool
	PoweredBy   string
	ProxyPath   string
}

func loadConfig(app *App) error {
	app.config = new(appConfig)
	data, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = json.Unmarshal(data, app.config)
	}
	return err
}
