package config

import (
	"log"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func New() {
	k := koanf.New(".")
	if err := k.Load(file.Provider("./config.yml"), yaml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	initVariables(k)

}
