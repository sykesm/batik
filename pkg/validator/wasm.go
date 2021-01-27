// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/pkg/errors"

	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
	"github.com/sykesm/batik/pkg/validator/utxo"
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
	validator, err := utxo.NewValidator(w.store, w.module)
	if err != nil {
		return nil, errors.WithMessage(err, "failed creating validator")
	}
	return validator.Validate(req)
}
