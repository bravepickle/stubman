package main

import (
	//	"bufio"
	"net/http"
	"net/http/httptest"
	//	"os"
	//	"regexp"
	//	"strings"
	"testing"
)

const testDbName = `test_db.sqlite`

func initCfgForTest(t *testing.T) {
	t.Log(`> Debug: disabled`)
	Debug = false

	Config.Stubman.Disabled = false
	Config.Stubman.Db.DbName = testDbName
}

func initTestDb(t *testing.T) {
	t.Log(`Init test DB:`, Config.Stubman.Db.DbName)
	_, err := NewDb(Config.Stubman.Db.DbName, true)
	if err != nil {
		t.Fatal(err)
	}
}

func TestServerStatic(t *testing.T) {
	t.Log(`>> Test running Server: check static files`)

	initCfgForTest(t)
	initTestDb(t)
	dbHelper := NewHelperDb(t)
	dbHelper.resetDbSchema()

	mux := http.NewServeMux()
	if err := initStubman(mux); err != nil {
		t.Fatal(err)
	}

	//	n := &http.Han{}

	srv := httptest.NewServer(mux)
	defer srv.Close()

	t.Log(`URL:`, srv.URL)

	testServerStatic(t, srv)
}

func testServerStatic(t *testing.T, srv *httptest.Server) {
	var testData = []struct {
		url         string
		code        int
		contentType string
	}{
		{
			url:         `/favicon.ico`,
			code:        200,
			contentType: `image/x-icon`,
		},
		{
			url:         `/unknown.html`,
			code:        404,
			contentType: `text/html; charset=utf-8`,
		},
	}

	var uri string
	for _, data := range testData {
		uri = srv.URL + data.url
		t.Log("\tChecking URI:", data.url)
		resp, err := http.Get(uri)
		if err != nil {
			t.Error("\t\tFailed to read", uri, `:`, err)
		}

		contentType := resp.Header.Get(`Content-Type`)
		if contentType != data.contentType {
			t.Errorf("\t\tFailed to assert MIME types. Expecting: \"%s\", received: \"%s\"",
				data.contentType, contentType)
		}

		if resp.StatusCode != data.code {
			t.Errorf("\t\tFailed to assert response status code. Expecting: %d, received: %d",
				data.code, resp.StatusCode)
		}

		//		time.Sleep(time.Second * 0.25) // wait for some time before checking - file data can be not flushed yet. Async write
		//		assertFileContains(t, data.url, Config.Log.Request.Output)

		//		if ok, _ := regexp.MatchString(`\.html`, data.url); ok { // check html extensions only
		//			assertFileContains(t, `Request `+data.url+` was received`, Config.Log.Response.Output)
		//		}
	}
}
