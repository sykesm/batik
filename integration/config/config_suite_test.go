// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

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
