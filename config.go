package goslim

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	// "fmt"
)

type config struct {
	Mode       string
	ListenAddr string
	ListenPort string
	PubDir     string
	Db         struct {
		SQLDriver string
		Protocol  string
		Hostname  string
		Port      string
		Username  string
		Password  string
		Dbname    string
		Charset   string
	}
}

func (c *config) loadConfig(configPath string) error {

	yamlS, readErr := ioutil.ReadFile(configPath)
	if readErr != nil {
		println(readErr.Error())

		return readErr
	}

	err := yaml.Unmarshal(yamlS, c)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%v", c)

	return nil
}
