package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	//	"strings"
	//	"time"
)

const defaultConfigPath = `./config.yaml`
const argCfgInit = `config:init`
const argDbInit = `db:init`
const argDbImport = `db:import`
const prefixPathStubman = `/stubman`

var optHelp bool
var cfgPath string
var Debug bool
var Config ConfigStruct

func init() {
	flag.BoolVar(&Debug, `debug`, false, `Enable debug mode`)
	flag.BoolVar(&optHelp, `help`, false, `Print command usage help`)
	flag.StringVar(&cfgPath, `f`, defaultConfigPath, `Path to config file in YAML format`)
}

func main() {
	flag.Parse()

	if optHelp {
		printAppUsage()
		return
	}

	if !initConfig(cfgPath, &Config) && flag.Arg(0) != argCfgInit {
		return
	}

	if !parseAppInput(cfgPath) {
		return
	}

	if Debug {
		fmt.Println(`Debug enabled`)
	}

	mux := http.NewServeMux()
	err := initStubman(mux)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	//favicon
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, `favicon.ico`)
	})

	mux.HandleFunc("/favicon.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, `favicon.png`)
	})

	prefixLen := len(prefixPathStubman) + 1
	mux.HandleFunc("/stubman/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[prefixLen:])
	})

	// handle the rest of URIs
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if !Config.Stubman.Disabled {
			//			ProxyRequest(w, req)
		} else {
			w.Header().Set(`X-Stubman-Page`, `true`)

			w.WriteHeader(http.StatusNotFound)

			page := PageError{Title: `404 Not Found`, Message: `Page "` + req.URL.String() + `" not found`}
			RenderErrorPage(`error.tpl`, page, w)
		}
	})

	if Debug {
		fmt.Printf("Listening to: %s\n", Config.App.String())
	}

	http.ListenAndServe(Config.App.String(), mux)
}

// init Stubman
func initStubman(mux *http.ServeMux) error {
	_, err := NewDb(Config.Stubman.Db.DbName, true)
	if err != nil {
		return err
	}

	AddStubmanCrudHandlers(prefixPathStubman, mux)

	// forward all static files to directory
	//	n.Use(negroni.NewStatic(http.Dir(``)))

	if Debug {
		fmt.Printf("Stubman path: http://%s%s/\n", Config.App.String(), prefixPathStubman)
	}

	return nil
}
