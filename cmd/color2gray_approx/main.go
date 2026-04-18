package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	cwa0 "github.com/takanoriyanagitani/go-color2gray/conv/wasm/wazero"
	"github.com/tetratelabs/wazero"
)

func sub(
	ctx context.Context,
	wasmLoc string,
	wasmByteMax int64,
	wasmLimitPage uint32,
	colorR, colorG, colorB float32,
) error {
	file, err := os.Open(wasmLoc) //nolint:gosec
	if nil != err {
		return err
	}
	defer file.Close() //nolint:errcheck

	limited := &io.LimitedReader{
		R: file,
		N: wasmByteMax,
	}

	wbytes, err := io.ReadAll(limited)
	if nil != err {
		return err
	}

	var rcfg wazero.RuntimeConfig = wazero.
		NewRuntimeConfig().
		WithMemoryLimitPages(wasmLimitPage)
	var mcfg wazero.ModuleConfig = wazero.NewModuleConfig()

	conv, err := cwa0.
		WasmBytes(wbytes).
		ToConverter(
			ctx,
			rcfg,
			mcfg,
		)
	if err != nil {
		return err
	}
	defer conv.Close(ctx) //nolint:errcheck

	converted, err := conv.ToGray(ctx, colorR, colorG, colorB)
	if nil != err {
		return err
	}

	fmt.Printf("gray: %v\n", converted) //nolint:forbidigo
	return nil
}

func main() {
	var wasmLoc string
	var wasmByteMax int64
	var wasmPageMax uint
	var colorR float64
	var colorG float64
	var colorB float64

	flag.StringVar(&wasmLoc, "wasm-path", "", "wasm path")
	flag.Int64Var(&wasmByteMax, "wasm-size-max", 131072, "wasm size max")
	flag.UintVar(&wasmPageMax, "wasm-page-max", 16, "wasm page max")
	flag.Float64Var(&colorR, "rfloat", 1.0, "red (0.0-1.0)")
	flag.Float64Var(&colorG, "gfloat", 0.50196078431372548, "green (0.0-1.0)")
	flag.Float64Var(&colorB, "bfloat", 0.0, "blue (0.0-1.0)")

	flag.Parse()

	if "" == wasmLoc {
		flag.Usage()
		os.Exit(1)
	}

	err := sub(
		context.Background(),
		wasmLoc,
		wasmByteMax,
		uint32(wasmPageMax), //nolint:gosec
		float32(colorR),
		float32(colorG),
		float32(colorB),
	)

	if nil != err {
		log.Printf("%v\n", err)
	}
}
