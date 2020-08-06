// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestEnvMap_Getenv(t *testing.T) {
	gt := NewGomegaWithT(t)

	envMap := EnvMap{
		"key": "value",
	}

	v, err := envMap.Getenv("key")
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(v).To(Equal("value"))

	v, err = envMap.Getenv("key2")
	gt.Expect(err).To(MatchError("$key2 is not defined"))
	gt.Expect(v).To(Equal(""))

}
