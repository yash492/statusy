package schema

import "embed"

//go:embed *.sql
var EmbedFS embed.FS
