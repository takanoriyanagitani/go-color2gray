#!/bin/sh

outname=./color2gray_approx_cli.wasm
mainpat=./cmd/color2gray_approx/main.go

GOOS=wasip1 GOARCH=wasm go \
	build \
	-o "${outname}" \
	-ldflags="-s -w" \
	"${mainpat}"
