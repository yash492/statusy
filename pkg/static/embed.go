package static

import "embed"

//go:embed services.yml
//go:embed config.yml
var Fs embed.FS
