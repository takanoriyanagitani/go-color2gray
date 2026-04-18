#!/bin/sh

bname="color2gray_approx"
bdir="./cmd/${bname}"
oname="${bdir}/${bname}"

mkdir -p "${bdir}"

go \
	build \
	-v \
	./...

go \
	build \
	-v \
	-o "${oname}" \
	"${bdir}"
