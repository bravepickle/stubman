# stubman
HTTP stub service written in Golang. Quick way to handle HTTP responses when target server is not ready yet. Contains WEB GUI for editing stubs

## Usage
Command usage

```
$ stubman -help

HTTP stub service.

Usage: stubman [options] [arg]

Options:
  -debug
    	Enable debug mode
  -f string
    	Path to config file in YAML format (default "./config.yaml")
  -help
    	Print command usage help

Arguments:
  config:init
    	initialize example config for running application. If file exists, then it will be reset to defaults
  db:init
    	initialize DB. If it exists, then DB will be reset
  db:import <file.sql>
    	import data from SQL file to DB. Second argument must be present with file path

Example:
  stubman config:init

```

## Features
- can work as stub service
- uses (SQLite)[http://github.com/mattn/go-sqlite3] as DB storage for stubman
- uses (YAML)[http://gopkg.in/yaml.v2] as config parser for stubman
- has WEB GUI for Stubs CRUD operations


## TODOs
- logger customizations: output source, naming, disabling etc.
- modes enable-disable
- allow redirect to another paths requests in proxy mode
- customize (change, add) request and response headers
- responses logger (disabled by default)
- separate file for error logging and access logs
- uniquely identify each request with responses, errors and keep those IDs in logs. Uniqueness can be partial
- cover with tests, benchmarks, profiling
- flag to disable stubman in whole and GUI side only
- proxy section: disable flag, 404 page by default
- switch stubman GUI themes 
- login use, credentials in config
- use SQLite as DB engine, add db:init command
- import/export stub data in GUI
- add to README examples how to setup configs for NGINX & Apache wit Midway App
- add multiplexer hendler: send copy of request to another address asynchoneously, without waiting for response (use goroutines?). Test it output
- add to config optional preconditions to start logging requests and responses, using RegEx
- in Stubman add button to generate CURL request for given stub
- log time taken for handling requests
- support method ANY which skips matching
- modify headers on save and keep as little as possible modifications on request handling
- add flag to strip GET params before matching given endpoint
- urldecode string before matching search pattern

## Usage Examples
### GUI Edit Page
![GUI Edit Page](https://monosnap.com/file/sJig2qUqVUfFFdydwSEq5gqvrBl0Be.png)

### GUI Listing Page
![GUI Listing Page](https://monosnap.com/file/W3y4wCF6rc2CGDPiA0hRP8pwEMr14K.png)

### Stub Example Results Page
![Stub Example Results Page](https://monosnap.com/file/tFjt6aXeJ1RtW6Jo82SP1vA9ySXlyp.png)
