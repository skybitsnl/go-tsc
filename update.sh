#!/usr/bin/env bash
set -euo pipefail

repo="microsoft/typescript"

version=${1:-}
if [ -z "$version" ]; then
	version=$(curl -fsSL "https://api.github.com/repos/$repo/releases/latest" |
		jq -r .tag_name | sed -n 's/^v\(.*\)$/\1/p')
fi
if [ -z "$version" ]; then
	echo "cannot determine the latest TypeScript release" >&2
	exit 1
fi

release="https://github.com/$repo/releases/download/v$version"

for platform in darwin-arm64 linux-arm64 linux-amd64; do
	os=$(echo "$platform" | awk -F- '{print $1}')
	arch=$(echo "$platform" | awk -F- '{print $2}')
	dlarch="$arch"
	if [ "$dlarch" = "amd64" ]; then
		dlarch="x64"
	fi
	curl -fsSL "$release/typescript-$os-$dlarch.tgz" | tar -xf - --strip-components 2 package/lib/tsc
	mv tsc tsc_${os}_${arch}
	gzip -9fn tsc_${os}_${arch}

	cat <<EOF >tsc_${os}_${arch}.go
package tsc

import _ "embed"

//go:embed tsc_${os}_${arch}.gz
var binary []byte
EOF
done

cat <<EOF >version.go
package tsc

// The version of TypeScript embedded by this module version.
const Version = "$version"
EOF

# Here, we assume the LICENSE, supporting files, and all .d.ts files
# in lib/ are equal among all OSes and architectures
rm -rf package
curl -fsSL "$release/typescript-linux-x64.tgz" | tar -xf -
rm package/lib/tsc
