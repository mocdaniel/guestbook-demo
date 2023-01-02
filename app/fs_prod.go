//go:build prod

package main

import (
	"embed"
	"io/fs"
)

//go:embed frontend/dist
var frontend embed.FS

func getFrontendAssets() (fs.FS, error) {
	return frontend, nil
}
