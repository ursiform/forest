// Copyright 2015 Afshin Darian. All rights reserved.
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

/*
Package forest_test contains tests and examples for package forest. The goal is
100% code coverage.
*/
package forest_test

import (
	"encoding/json"
	"github.com/ursiform/bear"
	"github.com/ursiform/forest"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func makeRequest(t *testing.T, app *forest.App,
	params *requested, want *wanted) *http.Response {
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
		t.Error(err)
		return nil
	}
	if response.Code != want.code {
		t.Errorf("%s %s want: %d (%s) got: %d %s, body: %s",
			method, path, want.code, http.StatusText(want.code), response.Code,
			http.StatusText(response.Code), string(responseBody))
		return nil
	}
	if responseData.Success != want.success {
		t.Errorf("%s %s should return success: %t", method, path, want.success)
		return nil
	}
	return &http.Response{Header: response.Header()}
}

func TestBasicOperation(t *testing.T) {
	debug := true
	path := "/foo"
	app := forest.New(debug)
	app.RegisterRoute(path, newRouter(app))
	params := &requested{method: "GET", path: path}
	want := &wanted{code: http.StatusOK, success: true}
	makeRequest(t, app, params, want)
}

func TestCookiesAndHeaders(t *testing.T) {
	debug := true
	cookieName := "foo"  // also in setCookie function of router
	cookieValue := "Foo" // also in setCookie function of router
	path := "/foo"
	app := forest.New(debug)
	app.RegisterRoute(path, newRouter(app))
	app.PoweredBy = "Testing-FTW"
	params := &requested{method: "GET", path: path}
	want := &wanted{code: http.StatusOK, success: true}
	response := makeRequest(t, app, params, want)
	if response == nil {
		return
	}
	if response.Header.Get("X-Powered-By") != app.PoweredBy {
		t.Errorf("app.PoweredBy header did not match response header: %s",
			response.Header.Get("X-Powered-By"))
	}
	for _, cookie := range response.Cookies() {
		if cookie.Name == cookieName && cookie.Value == cookieValue {
			return
		}
	}
	t.Errorf("cookie was not found")
}

func TestInitLog(t *testing.T) {
	// this test is just for complete coverage
	debug := false
	app := forest.New(debug)
	forest.InitLog(app, "", "foo message")           // undefined
	forest.InitLog(app, "initialize", "bar message") // install
	forest.InitLog(app, "install", "baz message")    // default
}

func TestInstallWare(t *testing.T) {
	debug := false
	app := forest.New(debug)
	message := "Foo handler installed"
	var handlerNil bear.HandlerFunc
	if err := app.InstallWare("Foo", handlerNil, message); err == nil {
		t.Errorf("app.InstallWare should reject nil handlers")
	}
	handler := func(http.ResponseWriter, *http.Request, *bear.Context) {}
	if err := app.InstallWare("Foo", handler, message); err != nil {
		t.Errorf("app.InstallWare failed: %s", err.Error())
	}
	// test duplicate ware installation
	if err := app.InstallWare("Foo", handler, message); err != nil {
		t.Errorf("app.InstallWare failed: %s", err.Error())
	}
}

func TestNonExistentWare(t *testing.T) {
	debug := false
	root := "/foo"
	path := "/foo/nonexistent"
	app := forest.New(debug)
	app.RegisterRoute(root, newRouter(app))
	params := &requested{method: "GET", path: path}
	want := &wanted{code: http.StatusInternalServerError, success: false}
	makeRequest(t, app, params, want)
}

func TestRetrievalDuration(t *testing.T) {
	debug := false
	durFoo := time.Hour * 1
	app := forest.New(debug)
	app.SetDuration("Foo", durFoo)
	if app.Duration("Foo") != durFoo {
		t.Errorf("SetDuration failed, want: %s got: %s",
			durFoo, app.Duration("Foo"))
	}
}

func TestRetrievalError(t *testing.T) {
	debug := false
	errFoo := "FOO_ERROR"
	app := forest.New(debug)
	app.SetError("Foo", errFoo)
	if app.Error("Foo") != errFoo {
		t.Errorf("SetError failed, want: %v got: %v", errFoo, app.Error("Foo"))
	}
}

func TestRetrievalMessage(t *testing.T) {
	debug := false
	msgFoo := "FOO_Message"
	app := forest.New(debug)
	app.SetMessage("Foo", msgFoo)
	if app.Message("Foo") != msgFoo {
		t.Errorf("SetMessage failed, want: %s got: %s",
			msgFoo, app.Message("Foo"))
	}
}

func TestServeFailure(t *testing.T) {
	debug := false
	path := "/foo"
	app := forest.New(debug)
	app.RegisterRoute(path, newRouter(app))
	go func() {
		if err := app.Serve(""); err == nil {
			t.Errorf("app.Serve should reject empty port string")
		}
	}()
}

func TestServeSuccess(t *testing.T) {
	debug := false
	path := "/foo"
	app := forest.New(debug)
	app.RegisterRoute(path, newRouter(app))
	go func() {
		if err := app.Serve(":31234"); err != nil {
			t.Errorf("app.Serve failed, %s", err.Error())
		}
	}()
}
