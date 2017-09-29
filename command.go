// operations related to command line interface input & parsing
package main

import (
	"flag"
	"fmt"
	"os"
)

const appConfigExample = `# common application settings
app:
  host: null   # any host will be checked
  port: 3001   # port to listen to
  base_path: . # base path to project templates, configs, etc.

stubman:
  disabled: false
  # DB settings
  db:
    dbname: ./data.sqlite
`

func printAppUsage() {
	fmt.Fprintln(os.Stderr, "HTTP stub service.\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [options] [arg]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, `Options:`)

	flag.PrintDefaults()

	fmt.Fprintf(os.Stderr, "\nArguments:\n")
	fmt.Fprintf(os.Stderr, "  %s\n    	initialize example config for running application. If file exists, then it will be reset to defaults\n", argCfgInit)
	fmt.Fprintf(os.Stderr, "  %s\n    	initialize DB. If it exists, then DB will be reset\n", argDbInit)
	fmt.Fprintf(os.Stderr, "  %s <file.sql>\n    	import data from SQL file to DB. Second argument must be present with file path\n", argDbImport)
	fmt.Fprintf(os.Stderr, "\nExample:\n  %s %s \n\n", os.Args[0], argCfgInit)
}

// parseAppInput parses input options and args from command line. Returns false when app should stop running
// after function execution
func parseAppInput(cfg string) bool {
	if flag.NArg() > 0 {
		switch flag.Arg(0) {
		case argCfgInit:
			if ok, err := saveToFile(appConfigExample, cfg); !ok || err != nil {
				fmt.Fprintf(os.Stderr, "Failed to init file \"%s\". Reason: %s\n", cfg, err.Error())
			} else {
				fmt.Printf("File \"%s\" was initialized successfully. Customize it and run application\n", cfg)
			}

			//			if err != nil {
			//				fmt.Fprintf(os.Stderr, "Failed to init file \"%s\". Reason: %s\n", cfg, err.Error())
			//			}

		case argDbInit:
			db, err := NewDb(Config.Stubman.Db.DbName, true)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to init DB. Reason: %s\n", err.Error())
			}

			if err = db.Reset(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to init DB. Reason: %s\n", err.Error())
			} else {
				fmt.Printf("DB \"%s\" was reset successfully.\n", Config.Stubman.Db.DbName)
			}

		case argDbImport:
			if flag.NArg() < 2 {
				fmt.Fprintln(os.Stderr, `Missing second argument with file path to import`)
				return false
			}

			db, err := NewDb(Config.Stubman.Db.DbName, true)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to init DB. Reason: %s\n", err.Error())
			}

			if err = db.ImportFromFile(flag.Arg(1)); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to import file \"%s\" to DB. Reason: %s\n", flag.Arg(1), err.Error())
			} else {
				fmt.Printf("File \"%s\" was successfully imported to %s.\n", flag.Arg(1), Config.Stubman.Db.DbName)
			}

		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", flag.Arg(0))
			printAppUsage()
		}

		return false
	}

	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File \"%s\" was not found. To init config run: %s %s\n", cfg, os.Args[0], argCfgInit)
		return false
	}

	return true
}
