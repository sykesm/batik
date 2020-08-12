// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
)

func TestBatik(t *testing.T) {
	gt := NewGomegaWithT(t)

	in := bytes.NewBuffer(nil)
	out := bytes.NewBuffer(nil)
	err := bytes.NewBuffer(nil)
	app := Batik(nil, ioutil.NopCloser(in), out, err)
	gt.Expect(app.Copyright).To(MatchRegexp("© Copyright IBM Corporation [\\d]{4}. All rights reserved."))

	gt.Expect(app.Flags).NotTo(BeEmpty())
	gt.Expect(app.Flags[0].Names()[0]).To(Equal("config"))
	gt.Expect(app.Commands).NotTo(BeEmpty())
	gt.Expect(app.Commands[0].Name).To(Equal("start"))
	gt.Expect(app.Commands[1].Name).To(Equal("status"))
}
