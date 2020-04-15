// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestBatik(t *testing.T) {
	gt := NewGomegaWithT(t)

	app := Batik(nil, nil, nil)
	gt.Expect(app.Copyright).To(MatchRegexp("© Copyright IBM Corporation [\\d]{4}. All rights reserved."))

	gt.Expect(app.Flags).To(BeEmpty())
	gt.Expect(app.Commands).To(BeEmpty())
}
