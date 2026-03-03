package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Debug      bool             `mapstructure:"debug"`
	PostgresDB PostgresDBConfig `mapstructure:"postgresdb"`
	// Interval in seconds
	ScrappingInterval int `mapstructure:"scrapping_interval"`
}

type PostgresDBConfig struct {
	ReadDB  DBConfig `mapstructure:"read"`
	WriteDB DBConfig `mapstructure:"write"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     uint   `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
}

func LoadConfig(filePath string) Config {
	path := "./config"
	if filePath != "" {
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
