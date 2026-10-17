// Package tsc embeds `microsoft/typescript` tsc releases into your Go binary,
// so you can perform TypeScript typechecking within your Go application. It
// also includes helper functions, like `RunCheck()`, to make this more
// convenient.
package tsc

import (
	"embed"
	"io/fs"
)

//go:embed package/LICENSE
var license string

//go:embed package/NOTICE.txt
var notice string

//go:embed package/lib/*.d.ts
var lib embed.FS

// Returns the license that was shipped along with the included tsc release.
func License() string {
	return license
}

// Returns the NOTICE file that was shipped along with the included tsc release.
func Notice() string {
	return notice
}

// Returns a filesystem of support files that were shipped along with the
// included tsc release. Currently, these are all .d.ts files. When running tsc,
// it expects these files to be present.
func Lib() fs.FS {
	f, err := fs.Sub(lib, "package/lib")
	if err != nil {
		panic(err)
	}
	return f
}
