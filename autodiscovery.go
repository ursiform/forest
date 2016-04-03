// Copyright 2016 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"fmt"
	"github.com/zeromq/gyre"
	"time"
)

const Network = "FOREST"

var (
	bus = make(chan []byte)
)

func autodiscoveryFailure(app *App, err error, index int) {
	message := fmt.Sprintf("autodiscovery failed[1]: %s", err.Error())
	InitLog(app, "warn", message)
}

func cluster(app *App) {
	defer close(bus)
	networkInterface := app.config.Autodiscovery.Interface
	name := app.config.Autodiscovery.Name
	node, err := gyre.New()
	if err != nil {
		autodiscoveryFailure(app, err, 1)
		return
	}
	if err = node.Start(); err != nil {
		autodiscoveryFailure(app, err, 2)
		return
	}
	defer node.Stop()
	if err = node.Join(Network); err != nil {
		autodiscoveryFailure(app, err, 3)
		return
	}
	if len(networkInterface) != 0 {
		if err = node.SetInterface(networkInterface); err != nil {
			autodiscoveryFailure(app, err, 4)
			return
		}
	} else {
		message := fmt.Sprintf(
			"%s.%s undefined in %s, using default",
			"autodiscovery", "interface", configFile)
		InitLog(app, "warn", message)
	}
	if len(name) != 0 {
		if err = node.SetName(name); err != nil {
			autodiscoveryFailure(app, err, 5)
			return
		}
	} else {
		message := fmt.Sprintf(
			"%s.%s undefined in %s, auto name: %s",
			"autodiscovery", "name", configFile, node.Name())
		InitLog(app, "initialize", message)
	}
	if err = node.SetHeader("service", app.config.Service.Name); err != nil {
		autodiscoveryFailure(app, err, 6)
		return
	}
	if err = node.SetHeader("version", app.config.Service.Version); err != nil {
		autodiscoveryFailure(app, err, 7)
		return
	}
	InitLog(app, "listen", "joined "+Network+" network")
	for {
		select {
		case event := <-node.Events():
			switch event.Type() {
			case gyre.EventEnter:
				if service, ok := event.Header("service"); ok {
					fmt.Printf("enter[1]: %s\n", service)
				} else {
					fmt.Printf("enter[1]\n")
				}
				fmt.Printf("enter[2]: %s\n", event.Name())
				if version, ok := event.Header("version"); ok {
					fmt.Printf("enter[3]: %s\n", version)
				} else {
					fmt.Printf("enter[3]\n")
				}
			case gyre.EventExit:
				fmt.Printf("exit: %s\n", event.Msg())
			case gyre.EventJoin:
				fmt.Printf("join: %s\n", event.Msg())
			case gyre.EventLeave:
				fmt.Printf("leave: %s\n", event.Msg())
			case gyre.EventShout:
				fmt.Printf("shout: %s\n", event.Msg())
			case gyre.EventWhisper:
				fmt.Printf("whisper: %s\n", event.Msg())
			}
		case message := <-bus:
			node.Shout(Network, message)
		}
	}
}

func loadNetworkDiscovery(app *App) {
	go cluster(app)
	return
}
