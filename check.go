package tsc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// A standard set of options to tsc which enables, among others, strict type
// checking against es2023.
func CompilerOptionsStrictES2023() map[string]any {
	return map[string]any{
		"strict":                     true,
		"noEmit":                     true,
		"noResolve":                  true,
		"target":                     "es2023",
		"lib":                        []string{"es2023"},
		"module":                     "esnext",
		"moduleResolution":           "bundler",
		"types":                      []string{},
		"skipLibCheck":               true,
		"noImplicitReturns":          true,
		"noFallthroughCasesInSwitch": true,
		"noUncheckedIndexedAccess":   true,
	}
}

// Like CompilerOptionsStrictES2023, but includes the "dom" library used
// in the browser.
func CompilerOptionsStrictES2023Browser() map[string]any {
	r := CompilerOptionsStrictES2023()
	r["lib"] = []string{"es2023", "dom"}
	return r
}

// RunCheck runs tsc with the given compilerOptions, then returns its error
// if type checking failed, or nil if it succeeded. You must pass bare filenames
// without subdirectories.
//
// It creates a temporary directory, uses Extract to extract tsc and all of its
// supporting files, then writes the files you're passing. It cleans up the
// temporary directory before returning. tsc is run using exec.CommandContext.
func RunCheck(ctx context.Context, compilerOptions map[string]any, sourceFiles map[string][]byte) error {
	dir, err := os.MkdirTemp("", "go-tsc-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	var fileNames []string
	for name, contents := range sourceFiles {
		if strings.Contains(name, "/") || name == "tsc" {
			return fmt.Errorf("invalid source file name: %q", name)
		}
		fileNames = append(fileNames, name)
		if err := os.WriteFile(filepath.Join(dir, name), contents, 0o600); err != nil {
			return err
		}
	}

	tsconfig, err := json.Marshal(map[string]any{
		"compilerOptions": compilerOptions,
		"files":           fileNames,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "tsconfig.json"), tsconfig, 0o600); err != nil {
		return err
	}

	tscPath, err := Extract(dir)
	if err != nil {
		return err
	}

	tscPath, err = filepath.Abs(filepath.Join(dir, tscPath))
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, tscPath, "-p", "tsconfig.json", "--pretty", "false")
	cmd.Dir = dir
	out := bytes.Buffer{}
	cmd.Stdout = &out
	cmd.Stderr = &out
	runErr := cmd.Run()
	if err := ctx.Err(); err != nil {
		return err
	}
	if runErr != nil {
		if out.String() == "" {
			return fmt.Errorf("tsc failed to run: %w", runErr)
		} else {
			return errors.New(out.String())
		}
	}

	return nil
}

// Extract extracts `tsc` for the running platform into the given directory,
// along with its supporting files. The directory must exist. This function
// returns where inside the directory `tsc` can be found. This is a relative
// path, so use `filepath.Abs(filepath.Join(dir, tscPath))` to find the absolute
// path.
func Extract(dir string) (tscPath string, err error) {
	fh, err := os.CreateTemp(dir, "tsc-*")
	if err != nil {
		return "", err
	}
	// This removes the file if we fail, and is a no-op once the rename succeeds
	defer os.Remove(fh.Name())

	if err := os.Chmod(fh.Name(), 0o700); err != nil {
		fh.Close()
		return "", err
	}
	if _, err := io.Copy(fh, BinaryReader()); err != nil {
		fh.Close()
		return "", err
	}
	if err := fh.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(fh.Name(), filepath.Join(dir, "tsc")); err != nil {
		return "", err
	}

	tscFiles, err := fs.ReadDir(Lib(), ".")
	if err != nil {
		return "", err
	}
	for _, tscFile := range tscFiles {
		b, err := fs.ReadFile(Lib(), tscFile.Name())
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(filepath.Join(dir, tscFile.Name()), b, 0o600); err != nil {
			return "", err
		}
	}

	return "tsc", nil
}
