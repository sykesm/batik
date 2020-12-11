// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utxo

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/pkg/errors"
	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
	"google.golang.org/protobuf/proto"
)

func modulesPath() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "..", "..", "..", "wasm", "modules")
}

type UTXOValidator struct {
	adapter *adapter
	engine  *wasmtime.Engine
	store   *wasmtime.Store

	modulePath string

	validateFunc *wasmtime.Func
}

func NewValidator() *UTXOValidator {
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	return &UTXOValidator{
		adapter:    &adapter{},
		engine:     engine,
		store:      store,
		modulePath: filepath.Join(modulesPath(), "utxotx", "target", "wasm32-unknown-unknown", "debug", "utxotx.wasm"),
	}
}

func (v *UTXOValidator) Init() error {
	if _, err := os.Stat(v.modulePath); os.IsNotExist(err) {
		return errors.Errorf("module does not exist at %s", v.modulePath)
	}

	module, err := wasmtime.NewModuleFromFile(v.engine, v.modulePath)
	if err != nil {
		return err
	}

	imports, err := v.newImports(module)
	if err != nil {
		return err
	}

	instance, err := wasmtime.NewInstance(
		v.store,
		module,
		imports,
	)
	if err != nil {
		return err
	}

	v.adapter.instance = instance
	v.adapter.memory = instance.GetExport("memory").Memory()
	v.validateFunc = instance.GetExport("validate").Func()

	return nil
}

func (v *UTXOValidator) Validate(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {

	resolved, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}

	v.adapter.resolved = resolved

	res, err := v.validateFunc.Call(99, len(resolved))
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
