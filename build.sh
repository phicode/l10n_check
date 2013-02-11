#!/bin/sh

build() {
	local out="bin/$3"
	env GOOS=$1 GOARCH=$2 go build -o "$out" github.com/PhiCode/l10n_check
	if [ $? -eq 0 ]; then
		echo "l10n_check for $1/$2 -> $out"
	fi
}

build linux   amd64 l10n_check
build windows amd64 l10n_check.exe


