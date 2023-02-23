package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type config struct {
	Db dbConfig `yaml:"db"`
}

var Env config

type dbConfig struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SslMode  string `yaml:"sslmode"`
	Dbname   string `yaml:"dbname"`
}

func Load() error {
	var config config

	bytes, err := os.ReadFile("../config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}
	Env = config

	fmt.Println(config)
	return nil
}
