package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Conf struct {
	Domain        string
	Port          string
	Debug         bool
	DbUser        string
	DbPassword    string
	DbName        string
	DbHost        string
	DbPort        string
	AllowedTokens []string
	Stats         bool
	Logs          bool
	LogsPath      string
	LogsLevel     int
}

func loadFromFile(file string) (*Conf, error) {
	if file == "" {
		file = "config.json"
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	var config = &Conf{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, errors.New("Error parsing config file: " + err.Error())
	}

	return config, nil
}

func newConf() *Conf {
	return &Conf{}
}
