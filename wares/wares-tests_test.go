// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

/*
Package wares_test contains tests and examples for package wares. The goal is
100% code coverage.
*/
package wares_test

import (
	"encoding/json"
	"github.com/ursiform/forest"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type requested struct {
	method string
	path   string
}

type wanted struct {
	code    int
	success bool
	data    interface{}
}

func makeRequest(t *testing.T, app *forest.App, params *requested, want *wanted) *http.Response {
	var request *http.Request
	method := params.method
	path := params.path
	request, _ = http.NewRequest(method, path, nil)
	response := httptest.NewRecorder()
	app.Router.ServeHTTP(response, request)
	responseData := new(forest.Response)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
		return nil
	}
	if err := json.Unmarshal(responseBody, responseData); err != nil {
		t.Errorf("unmarshal error: %v when attempting to read: %s", err, string(responseBody))
		return nil
	}
	if response.Code != want.code {
		t.Errorf("%s %s want: %d (%s) got: %d %s, body: %s", method, path,
			want.code, http.StatusText(want.code), response.Code, http.StatusText(response.Code), string(responseBody))
		return nil
	}
	if responseData.Success != want.success {
		t.Errorf("%s %s should return success: %t", method, path, want.success)
		return nil
	}
	return &http.Response{Header: response.Header()}
}

func TestAuthenticateFailure(t *testing.T) {
	debug := false
	method := "GET"
	root := "/foo"
	path := "/foo/authenticate/failure"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusUnauthorized, success: false}
	makeRequest(t, app, params, want)
}

func TestAuthenticateSuccess(t *testing.T) {
	debug := false
	method := "GET"
	root := "/foo"
	path := "/foo/authenticate/success"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestMethodNotAllowed(t *testing.T) {
	debug := false
	method := "OPTIONS"
	root := "/foo"
	path := "/foo"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusMethodNotAllowed, success: false}
	makeRequest(t, app, params, want)
}

func TestUnauthorized(t *testing.T) {
	debug := false
	method := "DELETE"
	root := "/foo"
	path := "/foo"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: method, path: path}
	want := &wanted{code: http.StatusUnauthorized, success: false}
	makeRequest(t, app, params, want)
}
