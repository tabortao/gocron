package embed

import (
	"embed"
	"io/fs"
)

//go:embed all:web/public
var files embed.FS

func StaticFS() (fs.FS, error) {
	return fs.Sub(files, "web/public")
}
