package config

import (
	"log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/yash492/statusy/pkg/static"
)

func New() {
	k := koanf.New(".")
	configBytes, err := static.Fs.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("error reading config yml: %v", err)
	}
	if err := k.Load(rawbytes.Provider(configBytes), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	initVariables(k)

}
