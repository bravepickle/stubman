// stub service entry point

package main

import (
	"bytes"
	"fmt"
	//	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const StaticPath = `static`
const viewsDir = `views`
const StubmanPathPrefix = `stubman`

type pathConcat struct {
	prefix string
}

// fullpath append prefix and path
func (p *pathConcat) fullPath(path string) string {
	buf := bytes.NewBufferString(p.prefix)
	buf.WriteString(path)

	return buf.String()
}

func init() {
	InitTemplates()
}

// AddGuiHandlers add all handlers for income requests that come to stub service
func AddStubmanCrudHandlers(prefix string, mux *http.ServeMux) {
	//	if Config.
	pcat := pathConcat{prefix}
	repo := NewStubRepo(nil)

	// list all stubs
	mux.HandleFunc(pcat.fullPath(`/`), func(w http.ResponseWriter, req *http.Request) {
		models, err := repo.FindAll()

		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)

			return
		}

		page := Page{HomePage: true, Data: models}
		RenderPage(`index.tpl`, page, w)
	})

	// create
	mux.HandleFunc(pcat.fullPath(`/create/`), func(w http.ResponseWriter, req *http.Request) {
		if req.Method == `POST` {
			req.ParseForm()
			stub := NewStubFromRequest(req)
			stub.Created = time.Now()
			stub.LastViewed = time.Now() // let's think that init of stub sets date to creation date

			id, err := repo.Insert(stub)
			if err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			w.Header().Add(`Location`, fmt.Sprintf(pcat.fullPath(`/edit/%d`), id))
			w.WriteHeader(http.StatusFound)

			return
		}

		model := NewNullObjectStub()
		page := Page{CreatePage: true, Data: model}
		RenderPage(`create.tpl`, page, w)
	})

	pathRegId := regexp.MustCompile(`\d+$`)
	// edit
	mux.HandleFunc(pcat.fullPath(`/edit/`), func(w http.ResponseWriter, req *http.Request) {
		id := pathRegId.FindString(req.URL.Path)
		if id == `` {
			w.Write([]byte(`Not Found`))
			w.WriteHeader(http.StatusNotFound)

			return
		}

		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		model, err := repo.Find(idNum)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if model.Id == 0 {
			w.Write([]byte(`Not Found`))
			w.WriteHeader(http.StatusNotFound)

			return
		}

		if req.Method == `POST` {
			req.ParseForm()
			stub := NewStubFromRequest(req)
			stub.Id = int64(idNum)

			err := repo.Update(stub)
			if err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			w.Header().Add(`Location`, fmt.Sprintf(pcat.fullPath(`/edit/%d`), idNum))
			w.WriteHeader(http.StatusFound)

			return
		}

		page := Page{EditPage: true, Data: model}
		RenderPage(`edit.tpl`, page, w)
	})

	// delete
	mux.HandleFunc(pcat.fullPath(`/delete/`), func(w http.ResponseWriter, req *http.Request) {
		id := pathRegId.FindString(req.URL.Path)
		if id == `` {
			w.Write([]byte(`Not Found`))
			w.WriteHeader(http.StatusNotFound)

			return
		}

		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)

			return
		}

		deleted, err := repo.Delete(idNum)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)

			return
		}

		if deleted {
			w.Write([]byte(`Not Found`))
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	})

	viewsStmt, err := repo.PrepareUpdateView()
	if err != nil {
		log.Fatalln(err.Error())
	}

	searchStmt, err := repo.PrepareSelectStubByRequest()
	if err != nil {
		log.Fatalln(err.Error())
	}

	// handle the rest of URIs
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if Config.Stubman.Disabled {
			w.Write([]byte(`Successfully received request: ` + req.Method + ` ` + req.URL.String()))
		} else {
			found, err := searchStmt.Query(req.Method, req.RequestURI)
			if err != nil {
				log.Println(req.RequestURI, `: `, err.Error())
				return
			}

			model := Stub{}
			for found.Next() {
				if err := found.Scan(&model.Id, &model.Name, &model.Response); err != nil {
					log.Println(req.RequestURI, `: `, err.Error())
					return
				}
				model.Decode()
			}

			if model.Id == 0 {
				w.Header().Set(`X-Stubman-Page`, `true`)
				w.WriteHeader(http.StatusNotFound)

				page := PageError{Title: `404 Not Found`, Message: `Page "` + req.URL.String() + `" not found`}
				RenderErrorPage(`error.tpl`, page, w)
			} else {
				for _, h := range model.ResponseParsed.Headers {
					arr := strings.SplitN(h, `:`, 2)
					if len(arr) == 2 {
						w.Header().Add(strings.TrimSpace(arr[0]), strings.TrimSpace(arr[1]))
					} else {
						fmt.Println(`Skipping broken header`, h, `for`, req.RequestURI)
					}
				}

				w.WriteHeader(model.ResponseParsed.StatusCode)
				w.Write([]byte(model.ResponseParsed.Body))

				go func(id int64) { viewsStmt.Exec(id) }(model.Id) // non-blocking mode for update views
			}
		}
	})
}

func NewStubFromRequest(req *http.Request) *Stub {
	stub := &Stub{}

	stub.Name = string(req.Form.Get(`stub_name`))
	stub.RequestMethod = string(req.Form.Get(`request_method`))
	stub.RequestUri = string(req.Form.Get(`request_uri`))
	stub.RequestParsed.Body = string(req.Form.Get(`request[body]`))
	for _, val := range req.Form[`request[headers][]`] {
		stub.RequestParsed.Headers = append(stub.RequestParsed.Headers, string(val))
	}

	stub.ResponseParsed.StatusCode, _ = strconv.Atoi(strings.TrimSpace(req.Form.Get(`resposne[status_code]`)))
	stub.ResponseParsed.Body = string(req.Form.Get(`response[body]`))
	for _, val := range req.Form[`response[headers][]`] {
		stub.ResponseParsed.Headers = append(stub.ResponseParsed.Headers, string(val))
	}

	stub.Encode()

	return stub
}
