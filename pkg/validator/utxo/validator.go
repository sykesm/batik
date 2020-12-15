// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utxo

import (
	"os"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/pkg/errors"
	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
	"google.golang.org/protobuf/proto"
)

// UTXOValidator implements the validator.Validator interface and provides
// a custom validator for UTXO transaction validation. Currently the web assembly
// module handles transaction signature verification.
type UTXOValidator struct {
	adapter *adapter
	engine  *wasmtime.Engine
	store   *wasmtime.Store
	module  *wasmtime.Module
}

func NewValidator(modulePath string) (*UTXOValidator, error) {
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		return nil, errors.Errorf("module does not exist at %s", modulePath)
	}

	module, err := wasmtime.NewModuleFromFile(engine, modulePath)
	if err != nil {
		return nil, err
	}

	return &UTXOValidator{
		adapter: &adapter{},
		engine:  engine,
		store:   store,
		module:  module,
	}, nil
}

func (v *UTXOValidator) Validate(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
	imports, err := v.newImports(v.module)
	if err != nil {
		return nil, err
	}

	instance, err := wasmtime.NewInstance(
		v.store,
		v.module,
		imports,
	)
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
