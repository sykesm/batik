// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package wasm

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bytecodealliance/wasmtime-go"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"

	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
)

func TestMain(m *testing.M) {
	cmd := exec.Command("cargo", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Join("modules", "utxotx")

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestBasic(t *testing.T) {
	gt := NewGomegaWithT(t)
	modulePath := filepath.Join("modules", "utxotx", "target", "wasm32-unknown-unknown", "debug", "utxotx.wasm")
	gt.Expect(modulePath).To(BeAnExistingFile())

	adapter := &adapter{}
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	module, err := wasmtime.NewModuleFromFile(engine, modulePath)
	gt.Expect(err).NotTo(HaveOccurred())

	var importedFuncs []*wasmtime.Extern
	for _, imp := range module.Imports() {
		var fn *wasmtime.Func
		switch imp.Module() {
		case "batik":
			switch imp.Name() {
			case "log":
				fn = wasmtime.WrapFunc(store, adapter.log)
			case "read":
				fn = wasmtime.WrapFunc(store, adapter.read)
			case "write":
				fn = wasmtime.WrapFunc(store, adapter.write)
			}
		}
		if fn == nil {
			panic(fmt.Sprintf("import %s::%s not found", imp.Module(), imp.Name()))
		}
		importedFuncs = append(importedFuncs, fn.AsExtern())
	}

	instance, err := wasmtime.NewInstance(
		store,
		module,
		importedFuncs,
	)
	gt.Expect(err).NotTo(HaveOccurred())

	adapter.instance = instance
	adapter.memory = instance.GetExport("memory").Memory()

	var validateRequest validationv1.ValidateRequest
	validateRequest.ResolvedTransaction = &validationv1.ResolvedTransaction{
		Txid: []byte("transaction-id"),
	}
	resolved, err := proto.Marshal(&validateRequest)
	gt.Expect(err).NotTo(HaveOccurred())
	adapter.resolved = resolved

	validate := instance.GetExport("validate").Func()
	res, err := validate.Call(99)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(res).To(BeEquivalentTo(0))

	var validateResponse validationv1.ValidateResponse
	err = proto.Unmarshal(adapter.response[4:], &validateResponse)
	gt.Expect(err).NotTo(HaveOccurred())
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
	var written int
	if a.idx < 4 {
		binary.BigEndian.PutUint32(buf, uint32(len(a.resolved)))
		written = 4
	} else {
		idx := a.idx - 4
		written = copy(buf, a.resolved[idx:idx+int(buflen)])
	}
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
