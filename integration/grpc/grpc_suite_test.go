// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/sykesm/batik/integration"
)

const testTimeout = 10 * time.Second

var binPaths paths

type paths struct {
	Batik            string
	WASMSigValidator string
}

var _ = SynchronizedBeforeSuite(func() []byte {
	batikPath, err := gexec.Build("github.com/sykesm/batik/cmd/batik")
	Expect(err).NotTo(HaveOccurred())

	cargoCmd := exec.Command("cargo", "build", "--target", "wasm32-unknown-unknown")
	cargoCmd.Dir = filepath.Join("..", "..", "wasm", "modules", "utxotx")

	wasmSigValidatorPath, err := filepath.Abs(
		filepath.Join(
			cargoCmd.Dir,
			"target",
			"wasm32-unknown-unknown",
			"debug",
			"utxotx.wasm",
		),
	)
	Expect(err).NotTo(HaveOccurred())

	cargoBuild, err := gexec.Start(cargoCmd, nil, nil)
	Expect(err).NotTo(HaveOccurred())
	Eventually(cargoBuild, time.Minute).Should(gexec.Exit(0))

	payload, err := json.Marshal(paths{
		Batik:            batikPath,
		WASMSigValidator: wasmSigValidatorPath,
	})
	Expect(err).NotTo(HaveOccurred())

	return []byte(payload)
}, func(payload []byte) {
	err := json.Unmarshal(payload, &binPaths)
	Expect(err).NotTo(HaveOccurred())
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	gexec.CleanupBuildArtifacts()
})

func StartPort() int {
	return integration.GRPCBasePort.StartPortForNode()
}

func TestGrpc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gRPC Suite")
}
