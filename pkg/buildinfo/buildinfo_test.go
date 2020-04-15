// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package buildinfo

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestFullVersion(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FullVersion()).To(Equal("dev-unknown"))
	})

	t.Run("Updated", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		defer func(version, commit string) {
			Version = version
			GitCommit = commit
		}(Version, GitCommit)

		Version = "new-version"
		GitCommit = "new-git-commit"
		gt.Expect(FullVersion()).To(Equal("new-version-new-git-commit"))

		GitCommit = ""
		gt.Expect(FullVersion()).To(Equal("new-version"))
	})
}

func TestBuilt(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		defer func() { now = time.Now }()

		now = func() time.Time {
			ts, err := time.Parse("2006-01-02", "2020-02-02")
			gt.Expect(err).NotTo(HaveOccurred())
			return ts
		}
		gt.Expect(Built()).To(Equal(now()))
	})

	t.Run("Updated", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		defer func(ts string) { Timestamp = ts }(Timestamp)

		Timestamp = "2020-02-02T00:00:00Z"
		ts, err := time.Parse("2006-01-02", "2020-02-02")
		gt.Expect(err).NotTo(HaveOccurred())

		gt.Expect(Built()).To(Equal(ts))
	})
}
