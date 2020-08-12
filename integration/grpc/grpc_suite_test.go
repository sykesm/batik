// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/sykesm/batik/integration"
)

const testTimeout = 1 * time.Second

var (
	batikPath string
)

var _ = SynchronizedBeforeSuite(func() []byte {
	batikPath, err := gexec.Build("github.com/sykesm/batik/cmd/batik")
	Expect(err).NotTo(HaveOccurred())

	payload, err := json.Marshal(batikPath)
	Expect(err).NotTo(HaveOccurred())

	return payload
}, func(payload []byte) {
	err := json.Unmarshal(payload, &batikPath)
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
	RunSpecs(t, "Grpc Suite")
}
