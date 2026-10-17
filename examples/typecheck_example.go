package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	tsc "github.com/skybitsnl/go-tsc/v7"
)

func main() {
	flag.Parse()

	fileNames := flag.Args()
	if len(fileNames) == 0 {
		log.Fatalf("Usage: %s file.ts [file.ts [...]]", os.Args[0])
	}

	files := map[string][]byte{}
	for _, name := range fileNames {
		body, err := os.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		files[filepath.Base(name)] = body
	}

	ctx := context.Background()
	err := tsc.RunCheck(ctx, tsc.CompilerOptionsStrictES2023Browser(), files)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("no errors")
	}
}
