// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
)

type WASM struct {
	store  *wasmtime.Store
	module *wasmtime.Module
}

func NewWASM(engine *wasmtime.Engine, asm []byte) (*WASM, error) {
	store := wasmtime.NewStore(engine)
	module, err := wasmtime.NewModule(engine, asm)
	if err != nil {
		return nil, err
	}
	return &WASM{store: store, module: module}, nil
}

func (w *WASM) Validate(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
	v := &UTXOValidator{
		adapter: &adapter{},
		store:   w.store,
		module:  w.module,
	}
	return v.Validate(req)
}

// UTXOValidator implements the validator.Validator interface and provides
// a custom validator for UTXO transaction validation. Currently the web assembly
// module handles transaction signature verification.
type UTXOValidator struct {
	adapter *adapter
	store   *wasmtime.Store
	module  *wasmtime.Module
}

func (v *UTXOValidator) Validate(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
	imports, err := v.newImports(v.module)
	if err != nil {
		return nil, err
	}

	instance, err := wasmtime.NewInstance(v.store, v.module, imports)
	if err != nil {
		return nil, err
	}

	v.adapter.instance = instance
	v.adapter.memory = instance.GetExport("memory").Memory()

	resolved, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	v.adapter.resolved = resolved

	res, err := instance.GetExport("validate").Func().Call(99, len(resolved))
	if err != nil {
		return nil, err
	}

	code, ok := res.(int32)
	if !ok {
		return nil, errors.Errorf("unrecognized return value: %v", res)
	}
	if code != 0 {
		return nil, errors.Errorf("validate failed, return code: %d", res)
	}

	validateResponse := &validationv1.ValidateResponse{}
	if err := proto.Unmarshal(v.adapter.response, validateResponse); err != nil {
		return nil, err
	}

	return validateResponse, nil
}

func (v *UTXOValidator) newImports(module *wasmtime.Module) ([]*wasmtime.Extern, error) {
	var importedFuncs []*wasmtime.Extern
	for _, imp := range module.Imports() {
		var fn *wasmtime.Func
		switch imp.Module() {
		case "batik":
			switch imp.Name() {
			case "log":
				fn = wasmtime.WrapFunc(v.store, v.adapter.log)
			case "read":
				fn = wasmtime.WrapFunc(v.store, v.adapter.read)
			case "write":
				fn = wasmtime.WrapFunc(v.store, v.adapter.write)
			}
		}
		if fn == nil {
			return nil, errors.Errorf("import %s::%s not found", imp.Module(), imp.Name())
		}
		importedFuncs = append(importedFuncs, fn.AsExtern())
	}

	return importedFuncs, nil
}

type adapter struct {
	instance *wasmtime.Instance
	memory   *wasmtime.Memory
	resolved []byte
	idx      int
	response []byte
}

func (a *adapter) read(streamID, addr, buflen int32) int32 {
	buf := a.memory.UnsafeData()[addr:]
	idx := a.idx
	written := copy(buf, a.resolved[idx:idx+int(buflen)])
	a.idx += written
	return int32(written)
}

func (a *adapter) write(streamID, addr, buflen int32) int32 {
	buf := a.memory.UnsafeData()[addr:]
	a.response = append(a.response, buf[:buflen]...)
	return buflen
}

func (a *adapter) log(buf, buflen int32) {
	fmt.Printf("%s\n", a.memory.UnsafeData()[buf:buf+buflen])
}
