// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var moddir = filepath.Join("..", "..", "wasm", "modules", "utxotx")

func TestMain(m *testing.M) {
	cmd := exec.Command("cargo", "build", "--target", "wasm32-unknown-unknown")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = moddir

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
