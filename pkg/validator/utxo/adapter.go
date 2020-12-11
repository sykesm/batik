// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utxo

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"
)

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
