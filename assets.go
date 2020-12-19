package manta

import (
	"embed"
)

const UIPrefix = `ui/build`

//go:embed ui/build
var Assets embed.FS
