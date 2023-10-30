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

	InitTemplates()

	if _, err := NewDb(Config.Stubman.Db.DbName, true); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	mux := http.NewServeMux()
	err := initStubman(mux)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if Debug {
		fmt.Printf("Listening to: %s\n", Config.App.String())
	}

	http.ListenAndServe(Config.App.String(), mux)
}

func initStaticHandlers(mux *http.ServeMux) {
	//favicon
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, Config.App.BasePath+`/favicon.ico`)
	})

	mux.HandleFunc("/favicon.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, Config.App.BasePath+`/favicon.png`)
	})

	prefixLen := len(prefixPathStubman) + 1
	mux.HandleFunc("/stubman/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, Config.App.BasePath+`/`+r.URL.Path[prefixLen:])
	})
}

// init Stubman handlers
func initStubman(mux *http.ServeMux) error {
	initStaticHandlers(mux)
	AddStubmanCrudHandlers(prefixPathStubman, mux)

	if Debug {
		baseUri := Config.App.BaseUri
		if baseUri == `` {
			baseUri = `/`
		}

		host := Config.App.Host
		if host == `` {
			host = `0.0.0.0`
		}

		if Config.App.Port != `80` {
			host += `:` + Config.App.Port
		}

		fmt.Printf("Stubman path: http://%s%s%s\nBase path: %s\n", host, Config.App.BaseUri, prefixPathStubman, Config.App.BasePath)
	}

	return nil
}
