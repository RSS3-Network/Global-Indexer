package docs

import "embed"

//go:embed openapi.yaml
var EmbedFS embed.FS
