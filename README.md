go-tsc
======

This module embeds `microsoft/typescript` tsc releases into your Go binary, so
you can perform TypeScript typechecking within your Go application.

But TypeScript is written in Go since 7.x, right? Why not simply import the
library instead of shipping the entire binary? Because `tsc` is implemented at
`/tsc/internal/` and Go disallows importing another package's internal
libraries. In order to call TypeScript from Go, we have to embed the binary.
This can only change if Microsoft moves the `tsc` implementation out of an
`internal` directory.

The simplest way to use this library is to use the utility function
`tsc.RunCheck` as follows:

```go
package main

import (
	"context"
	"log"

	"github.com/skybitsnl/go-tsc/v7"
)

var typescript = `
function formatCurrency(c: number) {
	console.log('$' + c);
}
formatCurrency("5");
`

func main() {
	ctx := context.Background()
	// Use tsc.CompilerOptionsStrictES2023() if not targeting the browser. When targeting
	// Node.js, pass a `.d.ts` file declaring the globals you use, such as `console`.
	err := tsc.RunCheck(ctx, tsc.CompilerOptionsStrictES2023Browser(), map[string][]byte{
		"target.ts": []byte(typescript),
	})
	log.Print(err)
}
```

Output:

```
target.ts(5,16): error TS2345: Argument of type 'string' is not assignable to
parameter of type 'number'.
```

Alternatively, you can use `tsc.Extract`, and write tsconfig.json and your TS
files yourself. Find the absolute path to tsc by
`filepath.Abs(filepath.Join(dir, tscPath))` and pass it to `exec.CommandContext`
with your own flags. Or, for full control, use `tsc.BinaryReader()` and
`tsc.Lib()` yourself.

The library also ships `tsc.License()` and `tsc.Notice()` from the TypeScript
distribution.

This library currently embeds the binaries for darwin-arm64, linux-arm64 and
linux-amd64. While building your application you'll download all of them; your
eventual binary will contain only the relevant binary.

Updating to a new TypeScript version
------------------------------------

- `./update.sh`
- Double check:
  - The `tsc` binaries are all changed
  - The `tsc` binaries belong to their respective platforms
  - The `package` directory contains only supporting files and .d.ts files
  - the Go tests succeed
- `git add` everything
- `git commit` mentioning the new TypeScript version
- `git tag v7.0.2` (replace with the new TypeScript version)
- `git push` and `git push --tags`

Making Go updates without bumping TypeScript
--------------------------------------------

Our Go library version number is intentionally linked to the TypeScript version
number. However, TypeScript versioning and the Go module versioning both follow
semver. If e.g. v7.0.4 of this module is already tagged, and we want to make
pure Go changes without bumping TypeScript, there is no way to express an update
to v7.0.4 without implying we ship TypeScript v7.0.5.

As such, the plan is to ship such Go updates in a prerelease, e.g. v7.0.5-fix.1.
That version would still contain TypeScript 7.0.4 but also any improvements on
the Go side. Such versions don't automatically ship to new users unless they opt
in explicitly. When TypeScript v7.0.5 then comes out, we publish our v7.0.5
which includes all the prerelease fixes.
