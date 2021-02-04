// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/sykesm/batik/integration"
)

const testTimeout = 10 * time.Second

var batikPath string

var _ = SynchronizedBeforeSuite(func() []byte {
	batikPath, err := gexec.Build("github.com/sykesm/batik/cmd/batik")
	Expect(err).NotTo(HaveOccurred())

	return []byte(batikPath)
}, func(payload []byte) {
	batikPath = string(payload)
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
