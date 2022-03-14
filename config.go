package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Conf struct {
	Domain               string
	Port                 string
	Debug                bool
	DbUser               string
	DbPassword           string
	DbName               string
	DbHost               string
	DbPort               string
	AllowedTokens        []string
	Stats                bool
	Logs                 bool
	LogsPath             string
	LogsLevel            int
	EnableFrontend       bool
	RateLimitPerIP       int
	RateLimitPerIPExpire int
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

func newConf() (*Conf, error) {
	configData, err := loadFromFile("config.json")
	if err != nil {
		return nil, err
	}
	return configData, nil
}

func loadConf() {
	configData, err := loadFromFile("config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config file: %s\n", err)
		fmt.Println("Error loading config file")
		l.Warning("Failed to load config file")
		os.Exit(1)
	}
	c = configData
	parseTokens(c.AllowedTokens)
	setRateLimit(c.RateLimitPerIP, c.RateLimitPerIPExpire)
	dbLink = c.DbUser + ":" + c.DbPassword + "@tcp(" + c.DbHost + ":" + c.DbPort + ")/" + c.DbName
}
