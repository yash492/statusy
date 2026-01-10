package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Debug      bool             `yaml:"debug"`
	PostgresDB PostgresDBConfig `yaml:"postgresdb"`
}

type PostgresDBConfig struct {
	ReadDB  DBConfig `yaml:"read"`
	WriteDB DBConfig `yaml:"read"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
}

func LoadConfig(filePath string) Config {
	path := "./config"
	if path != "" {
		path = filePath
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal the config: %s", err.Error())
	}

	return cfg
}

func (d DBConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", d.User, d.Password, d.Host, d.Port, d.Database)
}
