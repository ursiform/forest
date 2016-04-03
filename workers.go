// Copyright 2016 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

import (
	"fmt"
	"sync"
)

type worker struct {
	name    string
	version string
}

type workers struct {
	*sync.Mutex
	current int
	list    []worker
}

func (w *workers) add(name string, version string) int {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()
	for _, service := range w.list {
		if service.name == name {
			return len(w.list)
		}
	}
	fmt.Printf("\tadd [name: %s version: %s]\n", name, version)
	w.list = append(w.list, worker{name: name, version: version})
	return len(w.list)
}

func (w *workers) remove(name string) int {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()
	for index, service := range w.list {
		if service.name == name {
			fmt.Printf("\tremove [name: %s]\n", name)
			w.list = append(w.list[0:index], w.list[index+1:len(w.list)]...)
			break
		}
	}
	return len(w.list)
}

func newWorkers() *workers {
	w := &workers{}
	w.Mutex = new(sync.Mutex)
	w.list = make([]worker, 0)
	return w
}
