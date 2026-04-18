#!/bin/zsh

bt709approx=./grayapprox.wasm

wpath=./opt.wasm
wpath="${bt709approx}"

gwpath="/guest.d/read-only.d/grayapprox.wasm"

bin="./cmd/color2gray_approx/color2gray_approx"
wasi="./color2gray_approx_cli.wasm"

wazero \
  run \
  -mount "${PWD}:/guest.d/read-only.d:ro" \
  "${wasi}" \
  -wasm-page-max 1 \
  -wasm-size-max 1024 \
  -wasm-path "${gwpath}" \
  -rfloat 1.0 \
  -gfloat 0.50196078431372548 \
  -bfloat 0.0
