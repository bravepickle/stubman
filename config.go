// config-related operations & data
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// ConfigStruct struct contains main application config
type ConfigStruct struct {
	App     AppConfigStruct     `yaml:"app,omitempty"`
	Stubman StubmanConfigStruct `yaml:"stubman,omitempty"`
}

// StubmanConfigStruct contains Stubman settings
type StubmanConfigStruct struct {
	Disabled bool           `yaml:"disabled,omitempty"`
	Db       DbConfigStruct `yaml:"db,omitempty"`
}

// DbConfigStruct contains DB connection details for Stubman
type DbConfigStruct struct {
	DbName string `yaml:"dbname,omitempty"`
}

// HostPortString generates string scheme://host:port
func (t *DbConfigStruct) String() string {
	return fmt.Sprintf(`sqlite3://%s`, t.DbName)
}

// AppConfigStruct contain common application settings
type AppConfigStruct struct {
	Port     string
	Host     string
	BasePath string `yaml:"base_path"`
	BaseUri  string `yaml:"base_uri"`
}

func (t *AppConfigStruct) String() string {
	return fmt.Sprintf(`%s:%s`, t.Host, t.Port)
}

// initConfig parses config from file and puts it to config struct
func initConfig(cfgPath string, config *ConfigStruct) bool {
	if cfgPath == `` {
		cfgPath = defaultConfigPath
	}

	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open config: %s. Use %s %s to init config\n",
			cfgPath, os.Args[0], argCfgInit)
		return false
	}

	cfgFileString, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %s\n", err.Error())
		return false
	}

	err = yaml.Unmarshal(cfgFileString, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse config: %s\n", err.Error())
		return false
	}

	return true
}

func saveToFile(str string, cfgPath string) (bool, error) {
	file, err := os.Create(cfgPath)
	if err != nil {
		return false, err
	}

	_, err = file.WriteString(str)
	if err != nil {
		return false, err
	}

	return true, nil
}
