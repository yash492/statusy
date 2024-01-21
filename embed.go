package app

import "embed"

//go:embed services.yml
//go:embed models.sql
//go:embed config.yml
//go:embed _ui/build/*

var Fs embed.FS
