# stubman
HTTP stub service written in Golang. Quick way to handle HTTP responses when target server is not ready yet. Contains WEB GUI for editing stubs

## Usage
Command usage

```
$ midway -help

Web middleware app to log, proxy requests etc.

Usage: midway [options] [arg]

Options:
  -debug
    	Enable debug mode
  -f string
    	Path to config file in YAML format (default "./config.yaml")
  -help
    	Print command usage help
  -httptest.serve string
    	if non-empty, httptest.NewServer serves on this address and blocks

Arguments:
  config:init
    	initialize example config for running application. If file exists, then it will be reset to defaults
  db:init
    	initialize DB. If it exists, then DB will be reset
  db:import <file.sql>
    	import data from SQL file to DB. Second argument must be present with file path

Example:
  midway config:init

```

## Features
- can work as stub service
- uses (SQLite)[http://github.com/mattn/go-sqlite3] as DB storage for stubman
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

## Usage Examples
### GUI Add Page
![GUI Add Page](https://d1ro8r1rbfn3jf.cloudfront.net/ms_64293/t6tH2JisLEOnjy2GabVQJz1zcfpcdA/Stubman%2B%257C%2BStub%2BEdit%2B%25231%2B2017-01-30%2B10%2BAM10-13-40.png?Expires=1485850432&Signature=LQhq8FSePnYMh7cEDP316MGr0o~p5FKhLMB6qLp8pJyGRnvZGVK25M5nDRUxs1OGv70jnETfb-RUx34s0Zq3j82NfMb1uS3RQBPiyns163CwM1ZwQNFxU~-pVjq4g6lYhyw0ZCoGXvA2PdQJToTrYsKlVLTyrf9GcWa2x-2qrwKgZo1HGXbbj9RM52LCV2udjR22Dfe0-70f5BnYf2ZEDAoxnxSiqpYaWrRhsVhRnpMQHFQzgScIvcPsLV5LcigNAMOJKfaPE5PMOfhGUiWTOqkLuBHIvzCsWP0UM9eML0h61PGXL9u2RJWCKYlMExwcM3Jbpgj-~xJ40rNMmXH7dg__&Key-Pair-Id=APKAJHEJJBIZWFB73RSA)

### GUI Listing Page
![GUI Listing Page](https://d1ro8r1rbfn3jf.cloudfront.net/ms_64293/WwkjX6BEW9dQi5kjRdjwIojhRq7Fyk/Stubman%2B%257C%2BStubs%2BList%2B2017-01-30%2B10%2BAM10-15-14.png?Expires=1485850701&Signature=DTRVdaj0APBbYy6xLe3h7uwmh6ON61hK7z~1Mw-4bsPaB8~dks9gW4bmRTa9b6IDuTbIlOCcZxmItF-TKWKVN0FVZOh-B6uoXSdIqrgG1QMH~EKcF-0WatZPrLa6rXopT3FvRbscXVTc8YeIoCl16NDf4rnoA7NCc9MFwxjOIZlpgqDJATOm-KlPZKoK1fde06ebSSNKFWeYSa0k0QrueRwcGqGowhkz4s9tCl9sFOk1MyckCxRh514xUiC45-wIjAVSqzh8xirv7Hlzhq5tnefjFxwQYFDoyIgXHPHEXhmJYJlUs-2938qI1ZPzP54gNH~XnjFHZvRJLi6NOcnR5w__&Key-Pair-Id=APKAJHEJJBIZWFB73RSA)

### Stub Example Results Page
![Stub Example Results Page](https://d1ro8r1rbfn3jf.cloudfront.net/ms_64293/1hnfjFAgFAQfIKkCCdE84jaoX05x68/Hello%252C%2BWorld%2521%2B2017-01-30%2B10%2BAM10-14-27.png?Expires=1485850481&Signature=dPiftD~7zESmSQ7QM0r61fv4U7KZ7eHXhNcAIP3XURogMrHautBpHxulbDr43~8v7Ruxj5QKVFp4EbvQK9ynqKcXCeO1Teo4xcqsu1dNtpxIhXtVQq0qREl874dQ0X6-UVTZp4C-0W9OxJxRcsrnknrNS~NZwgTLckimCJ-ufbi~4f2xZXkVkzDbl0c93RKjJikzTcHLbe33V~gVq8O1axf8oYIuz4mTBFzE3hKiiFeafaE6pizgxNUL2mMsT~xDSH-acCzAytumaTe0oQGZ12uIt0zRkdB6ctGh-iZ84AY29WBCbiibimSj~wLH9Z2cHH-yVrzHvynZDdztjmqGnw__&Key-Pair-Id=APKAJHEJJBIZWFB73RSA)
