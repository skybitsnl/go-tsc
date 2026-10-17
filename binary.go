package tsc

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Returns an io.Reader which contains the tsc binary for your running
// platform. (The binary is uncompressed as you read.)
func BinaryReader() io.Reader {
	r, err := gzip.NewReader(bytes.NewReader(binary))
	if err != nil {
		panic(err)
	}
	return r
}
