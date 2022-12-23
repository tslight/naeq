package books

import "embed"

//go:embed *.json
var EFS embed.FS
