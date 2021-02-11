// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"os"
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

var (
	batikPath            string
	wasmSigValidatorPath string
)

const sigvalPath = "../../rust/sigval/target/wasm32-unknown-unknown/release/sigval.wasm"

var _ = SynchronizedBeforeSuite(func() []byte {
	_, err := os.Stat(sigvalPath)
	if os.IsNotExist(err) {
		cmd := exec.Command("make", "cargo-build")
		cmd.Dir = filepath.Join("..", "..")
		sess, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(sess).Should(gexec.Exit(0))
	}

	batikPath, err := gexec.Build("github.com/sykesm/batik/cmd/batik")
	Expect(err).NotTo(HaveOccurred())

	return []byte(batikPath)
}, func(payload []byte) {
	path, err := filepath.Abs(sigvalPath)
	Expect(err).NotTo(HaveOccurred())
	Expect(path).To(BeAnExistingFile())

	batikPath = string(payload)
	wasmSigValidatorPath = path
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
