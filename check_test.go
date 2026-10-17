package tsc_test

import (
	"context"
	"testing"

	"github.com/skybitsnl/go-tsc/v7"
)

func TestOk(t *testing.T) {
	files := map[string][]byte{
		"target.ts": []byte(`
function formatCurrency(c: number) {
  console.log('$' + c);
}
formatCurrency(5);
`),
	}

	ctx := context.Background()
	err := tsc.RunCheck(ctx, tsc.CompilerOptionsStrictES2023Browser(), files)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMissingDomError(t *testing.T) {
	files := map[string][]byte{
		"target.ts": []byte(`
function formatCurrency(c: number) {
  console.log('$' + c);
}
formatCurrency(5);
`),
	}

	ctx := context.Background()
	err := tsc.RunCheck(ctx, tsc.CompilerOptionsStrictES2023(), files)
	if err == nil {
		t.Fatal("error expected, but got nil")
	}
}

func TestError(t *testing.T) {
	files := map[string][]byte{
		"target.ts": []byte(`
function formatCurrency(c: number) {
  console.log('$' + c);
}
formatCurrency("5");
`),
	}

	ctx := context.Background()
	err := tsc.RunCheck(ctx, tsc.CompilerOptionsStrictES2023Browser(), files)
	if err == nil {
		t.Fatal("error expected, but got nil")
	}
}
