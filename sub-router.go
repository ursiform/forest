// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package forest

/*
SubRouter is the basic building block of forest applications. In the main.go
file where a service is initialized, app.RegisterRoute is called with a path
string and a SubRouter instance in order to logically group all of the endpoints
a particular service will answer. Since each App inherits from bear.Mux , it is
not strictly necessary to use SubRouter at all, it simply exists to provide a
convenient way to instantiate many endpoints in one place.
*/
type SubRouter interface {
	Route(path string)
}
