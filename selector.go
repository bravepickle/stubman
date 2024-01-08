// select which stub to use for given request
package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

// ByBestRequestMatch sort stubs by best matches to given request
// type []Stub
type ByBestRequestMatch struct {
	stubs   []Stub
	request *http.Request
}

func (a ByBestRequestMatch) Len() int      { return len(a.stubs) }
func (a ByBestRequestMatch) Swap(i, j int) { a.stubs[i], a.stubs[j] = a.stubs[j], a.stubs[i] }
func (a ByBestRequestMatch) Less(i, j int) bool {
	// TODO: sorting here

	// more loose match by URI
	if len(a.stubs[i].RequestUri) > len(a.stubs[j].RequestUri) {
		return false
	}

	// give priority to stub with more defined headers
	if len(a.stubs[i].RequestParsed.Headers) > len(a.stubs[j].RequestParsed.Headers) {
		return true
	}

	// more precise method
	if a.stubs[j].RequestMethod == `ANY` && a.stubs[i].RequestMethod != `ANY` {
		return true
	}

	// specified body
	if a.stubs[j].RequestParsed.Body != `` && a.stubs[i].RequestParsed.Body == `` {
		return true
	}

	// give priority to newer stub otherwise
	return a.stubs[j].Id < a.stubs[i].Id
}

func SortByRequest(req *http.Request, stubs *[]Stub) {
	sorted := ByBestRequestMatch{stubs: *stubs, request: req}
	sort.Sort(sorted)
}

// selectStub selects which Stub to use based on request
func selectStub(req *http.Request, searchStmt *sql.Stmt) (selected *Stub, err error) {
	found, err := searchStmt.Query(req.Method, req.RequestURI+`%`)

	log.Println(`[INFO] URI SEARCH: `, req.Method, req.RequestURI)

	if err != nil {
		log.Println(req.RequestURI, `: `, err.Error())

		return nil, err
	}

	var bodyStr string

	if req.Method != `GET` {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(`[ERROR] Failed reading body: `, err.Error(), `BODY:`, body)

			return nil, err
		}

		bodyStr = string(body)
		if bodyStr != `` {
			log.Println(`[INFO] REQUEST BODY: `, bodyStr)
		} else if req.Method != `GET` {
			log.Println(`[INFO] NO BODY`)
		}
	} else {
		bodyStr = `` // no body for GET method is possible
	}

	var stubs []Stub
	for found.Next() {
		model := Stub{}
		if err := found.Scan(&model.Id, &model.Name,
			&model.RequestUri, &model.RequestMethod, &model.Response, &model.Request); err != nil {
			log.Println(req.RequestURI, `: `, err.Error())

			return nil, err
		}

		model.Decode()

		if matchStubData(&model, req, bodyStr) {
			stubs = append(stubs, model)
		}
	}

	stubsNum := len(stubs)
	if stubsNum > 1 {
		SortByRequest(req, &stubs)

		log.Println(`[INFO] FOUND: `, stubs[0].Id, req.Method, req.RequestURI)

		return &stubs[0], nil
	} else if stubsNum == 1 {
		return &stubs[0], nil
	}

	log.Println(`[INFO] STUB NOT FOUND`)

	return selected, nil
}

func matchStubData(model *Stub, req *http.Request, body string) bool {
	// check if headers are set
	if len(model.RequestParsed.Headers) > 0 {
		for _, h := range model.RequestParsed.Headers {
			arr := strings.SplitN(h, `:`, 2)
			if !containsHeader(strings.TrimSpace(arr[0]), strings.TrimSpace(arr[1]), req.Header) {
				return false
			}
		}
	}

	if req.Method != `GET` && req.Method != `ANY` && model.RequestParsed.Body != `` {
		if body != model.RequestParsed.Body {
			return false
		}
	}

	return true
}
