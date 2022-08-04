package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/jackpf/kraken-scheduler/src/main/config/model"
)

func ParseConfigFile(fileName string) (*model.Config, error) {
	data, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	return ParseConfig(data)
}

func ParseConfig(data []byte) (*model.Config, error) {
	var config model.Config

	err := json.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
