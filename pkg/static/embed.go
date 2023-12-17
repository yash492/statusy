package static

import "embed"

//go:embed services.yml
//go:embed config.yml
//go:embed models.sql
var Fs embed.FS
