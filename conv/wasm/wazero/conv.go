package conv

import (
	"context"
	"errors"

	"github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"
)

var (
	ErrNilFunc error = errors.New("nil function")

	ErrInvalidResults error = errors.New("unexpected num of results")
)

type WasmFn struct{ wa.Function }

func (f WasmFn) ConvertRgbFloat(
	ctx context.Context,
	r, g, b float32,
) (float32, error) {
	var er uint64 = wa.EncodeF32(r)
	var eg uint64 = wa.EncodeF32(g)
	var eb uint64 = wa.EncodeF32(b)

	results, err := f.Function.Call(ctx, er, eg, eb)

	if nil != err {
		return 0.0, err
	}

	if 1 != len(results) {
		return 0.0, ErrInvalidResults
	}

	var result uint64 = results[0]
	var decoded float32 = wa.DecodeF32(result)

	return decoded, nil
}

type WasmMod struct{ wa.Module }

func (m WasmMod) Close(ctx context.Context) error {
	if nil == m.Module {
		return nil
	}
	return m.Module.Close(ctx)
}

func (m WasmMod) GetFunction(name string) (WasmFn, error) {
	var fnc wa.Function = m.Module.ExportedFunction(name)
	if nil == fnc {
		return WasmFn{}, ErrNilFunc
	}
	return WasmFn{Function: fnc}, nil
}

func (m WasmMod) GetDetector64() (WasmFn, error) {
	return m.GetFunction("luminance32f")
}

type Compiled struct{ wazero.CompiledModule }

func (c Compiled) Close(ctx context.Context) error {
	if nil == c.CompiledModule {
		return nil
	}
	return c.CompiledModule.Close(ctx)
}

type WasmRuntime struct{ wazero.Runtime }

func (r WasmRuntime) Close(ctx context.Context) error {
	if nil == r.Runtime {
		return nil
	}
	return r.Runtime.Close(ctx)
}

func (r WasmRuntime) Compile(
	ctx context.Context,
	wasm []byte,
) (Compiled, error) {
	cmod, err := r.Runtime.CompileModule(ctx, wasm)
	return Compiled{CompiledModule: cmod}, err
}

func (r WasmRuntime) Instantiate(
	ctx context.Context,
	compiled Compiled,
	cfg wazero.ModuleConfig,
) (WasmMod, error) {
	amod, err := r.Runtime.InstantiateModule(
		ctx,
		compiled.CompiledModule,
		cfg,
	)

	return WasmMod{Module: amod}, err
}

type WasmConfig struct{ wazero.RuntimeConfig }

type Converter struct {
	WasmRuntime
	Compiled
	WasmMod
	WasmFn
}

func (c Converter) Close(ctx context.Context) error {
	return errors.Join(
		c.WasmMod.Close(ctx),
		c.Compiled.Close(ctx),
		c.WasmRuntime.Close(ctx),
	)
}

func (c Converter) ToGray(
	ctx context.Context,
	r float32,
	g float32,
	b float32,
) (float32, error) {
	return c.WasmFn.ConvertRgbFloat(ctx, r, g, b)
}

type WasmBytes []byte

func (b WasmBytes) ToConverter(
	ctx context.Context,
	rcfg wazero.RuntimeConfig,
	mcfg wazero.ModuleConfig,
) (Converter, error) {
	var rtm wazero.Runtime = wazero.NewRuntimeWithConfig(
		ctx,
		rcfg,
	)
	var detector Converter
	detector.WasmRuntime = WasmRuntime{Runtime: rtm}

	compiled, err := rtm.CompileModule(ctx, b)
	if nil != err {
		return detector, err
	}
	detector.Compiled = Compiled{CompiledModule: compiled}

	instance, err := rtm.InstantiateModule(
		ctx,
		detector.Compiled.CompiledModule,
		mcfg,
	)
	if nil != err {
		return detector, err
	}
	detector.WasmMod = WasmMod{Module: instance}

	det64, err := detector.WasmMod.GetDetector64()
	if nil != err {
		return detector, err
	}
	detector.WasmFn = det64

	return detector, nil
}
