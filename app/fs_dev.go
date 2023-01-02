//go:build !prod

package main

import (
	"io/fs"
	"os"
)

func getFrontendAssets() (fs.FS, error) {
	return os.DirFS("./app"), nil
}
