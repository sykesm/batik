// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/gomega"
)

func TestSearchPath(t *testing.T) {
	gt := NewGomegaWithT(t)

	wd, err := os.Getwd()
	gt.Expect(err).NotTo(HaveOccurred())
	confDir, err := os.UserConfigDir()
	gt.Expect(err).NotTo(HaveOccurred())

	paths, err := SearchPath("franklin")
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(paths[0]).To(Equal(wd))
	gt.Expect(paths[1]).To(Equal(filepath.Join(confDir, "franklin")))
	if runtime.GOOS == "darwin" {
		gt.Expect(paths).To(HaveLen(3))
		gt.Expect(paths[2]).To(Equal(filepath.Join(os.Getenv("HOME"), ".config", "franklin")))
	}
}
