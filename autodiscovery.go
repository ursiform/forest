// Copyright 2016 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"fmt"
	"github.com/zeromq/gyre"
)

var (
	requests  = make(chan []byte)
	services  = make(map[string]*workers)
	directory = make(map[string]string)
)

func discoveryFailure(app *App, err error, code string) {
	message := fmt.Sprintf("[autodiscovery error %s]: %s", code, err.Error())
	InitLog(app, "warn", message)
}

func discover(app *App) {
	defer close(requests)
	node, err := gyre.New()
	if err != nil {
		discoveryFailure(app, err, errorDiscoveryInitialize)
		return
	}
	if app.config.Autodiscovery.Port > 0 {
		if err = node.SetPort(app.config.Autodiscovery.Port); err != nil {
			discoveryFailure(app, err, errorDiscoverySetPort)
			return
		}
	}
	networkInterface := app.config.Autodiscovery.Interface
	if len(networkInterface) != 0 {
		if err = node.SetInterface(networkInterface); err != nil {
			discoveryFailure(app, err, errorDiscoveryNetworkInterface)
			return
		}
	} else {
		message := fmt.Sprintf(
			"%s.%s not defined in %s, using default",
			"autodiscovery", "interface", configFile)
		InitLog(app, "warn", message)
	}
	name := app.config.Autodiscovery.Name
	if len(name) > 0 {
		if err = node.SetName(name); err != nil {
			discoveryFailure(app, err, errorDiscoverySetName)
			return
		}
	} else {
		message := fmt.Sprintf(
			"%s.%s not defined in %s, using name: %s",
			"autodiscovery", "name", configFile, node.Name())
		InitLog(app, "initialize", message)
	}
	if err = node.Start(); err != nil {
		discoveryFailure(app, err, errorDiscoveryStart)
		return
	}
	defer node.Stop()
	if err = node.Join(network); err != nil {
		discoveryFailure(app, err, errorDiscoveryJoin)
		return
	}
	if err = node.SetHeader("service", app.config.Service.Name); err != nil {
		discoveryFailure(app, err, errorDiscoveryServiceHeader)
		return
	}
	version := app.config.Service.Version
	if len(version) == 0 {
		version = "unknown"
	}
	if err = node.SetHeader("version", version); err != nil {
		discoveryFailure(app, err, errorDiscoveryVersionHeader)
		return
	}
	if port := app.config.Autodiscovery.Port; port > 0 {
		InitLog(app, "listen", fmt.Sprintf("%s:%d", network, port))
	} else {
		InitLog(app, "listen", fmt.Sprintf("autodiscovery [%s:%d]", network, 5670))
	}
	for {
		select {
		case event := <-node.Events():
			switch event.Type() {
			case gyre.EventEnter:
				fmt.Printf("[autodiscovery enter] %s\n", event.Name())
				service, ok := event.Header("service")
				version, _ := event.Header("version")
				name := event.Name()
				if ok {
					directory[name] = service
					if services[service] == nil {
						services[service] = newWorkers()
					}
					services[service].add(name, version)
				}
			case gyre.EventExit:
				fmt.Printf("[autodiscovery exit] %s\n", event.Name())
				name := event.Name()
				service, ok := directory[name]
				if ok {
					remaining := services[service].remove(name)
					if remaining == 0 {
						delete(services, service)
					}
					delete(directory, name)
				}
			case gyre.EventJoin:
				fmt.Printf("[autodiscovery join] %s\n", event.Name())
			case gyre.EventLeave:
				fmt.Printf("[autodiscovery leave] %s\n", event.Name())
				service, ok := event.Header("service")
				version, _ := event.Header("version")
				if ok {
					if services[event.Name()] == nil {
						services[event.Name()] = newWorkers()
						services[event.Name()].add(service, version)
					}
				}
			case gyre.EventShout:
				fmt.Printf("[autodiscovery shout] %s\n", event.Msg())
			case gyre.EventWhisper:
				fmt.Printf("[autodiscovery whisper] %s\n", event.Msg())
			}
		case request := <-requests:
			node.Shout(network, request)
		}
	}
}

func loadNetworkDiscovery(app *App) {
	if app.config.Autodiscovery.Silent {
		return
	}
	if len(app.config.Service.Name) == 0 {
		err := fmt.Errorf("%s.%s not defined in %s", "service", "name", configFile)
		discoveryFailure(app, err, errorDiscoveryServiceUndefined)
		return
	}
	go discover(app)
}
