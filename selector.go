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
//type []Stub
type ByBestRequestMatch struct {
	stubs   []Stub
	request *http.Request
}

func (a ByBestRequestMatch) Len() int      { return len(a.stubs) }
func (a ByBestRequestMatch) Swap(i, j int) { a.stubs[i], a.stubs[j] = a.stubs[j], a.stubs[i] }
func (a ByBestRequestMatch) Less(i, j int) bool {
	// TODO: sorting here
	return a.stubs[i].Id < a.stubs[j].Id
}

func SortByRequest(req *http.Request, stubs *[]Stub) {
	sorted := ByBestRequestMatch{stubs: *stubs, request: req}
	sort.Sort(sorted)

	log.Println(` -------- SORTED:`, sorted.stubs)
}

//// By is the type of a "less" function that defines the ordering of its Planet arguments.
//type By func(p1, p2 *Planet) bool

//// Sort is a method on the function type, By, that sorts the argument slice according to the function.
//func (by By) Sort(planets []Planet) {
//	ps := &planetSorter{
//		planets: planets,
//		by:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
//	}
//	sort.Sort(ps)
//}

//// planetSorter joins a By function and a slice of Planets to be sorted.
//type planetSorter struct {
//	planets []Planet
//	by      func(p1, p2 *Planet) bool // Closure used in the Less method.
//}

//// Len is part of sort.Interface.
//func (s *planetSorter) Len() int {
//	return len(s.planets)
//}

//// Swap is part of sort.Interface.
//func (s *planetSorter) Swap(i, j int) {
//	s.planets[i], s.planets[j] = s.planets[j], s.planets[i]
//}

//// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
//func (s *planetSorter) Less(i, j int) bool {
//	return s.by(&s.planets[i], &s.planets[j])
//}

// selectStub selects which Stub to use based on request
func selectStub(req *http.Request, searchStmt *sql.Stmt) (selected *Stub, err error) {
	found, err := searchStmt.Query(req.Method, req.RequestURI)
	if err != nil {
		log.Println(req.RequestURI, `: `, err.Error())

		return nil, err
	}

	var stubs []Stub
	for found.Next() {
		model := Stub{}
		if err := found.Scan(&model.Id, &model.Name, &model.Response, &model.Request); err != nil {
			log.Println(req.RequestURI, `: `, err.Error())

			return nil, err
		}

		model.Decode()

		if matchStubData(&model, req) {
			stubs = append(stubs, model)
			log.Println(`---MATCHED_STUB---`, model.Id, model.Name)
		}
	}

	if len(stubs) > 0 {
		log.Println(`---PRE_THE_STUBS---`, stubs)

		SortByRequest(req, &stubs)

		log.Println(`---THE_STUBS---`, stubs)

		return &stubs[0], nil
	}

	// TODO: select the best match from existing - more precise

	return selected, nil
}

func matchStubData(model *Stub, req *http.Request) bool {
	// check if headers are set
	if len(model.RequestParsed.Headers) > 0 {
		for _, h := range model.RequestParsed.Headers {
			arr := strings.SplitN(h, `:`, 2)
			if !containsHeader(strings.TrimSpace(arr[0]), strings.TrimSpace(arr[1]), req.Header) {

				log.Println(`=__= no header`, model.Id, ` = `, h)

				return false
			}
		}
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(req.RequestURI, `:`, err.Error())

		return false
	}

	if req.Method != `GET` && req.Method != `ANY` && model.RequestParsed.Body != `` {
		if string(body) != model.RequestParsed.Body {
			return false
		}
	}

	return true
}
