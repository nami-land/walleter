package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Database struct {
		Addr     string `yaml:"addr"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
}

func NewDevConfig() *Config {
	config, err := GetConfig("../config.dev.yaml")
	if err != nil {
		return nil
	}
	return config
}

func GetConfig(filePath string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
