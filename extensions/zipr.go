//  zipr.go -- Zip file reader extension for Goaldi

package extensions

import (
	"archive/zip"
	"goaldi"
)

func init() {
	goaldi.GoLib(zip.OpenReader, "zipreader", "name", "open a Zip file")
}
