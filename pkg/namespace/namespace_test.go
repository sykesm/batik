// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/tested"
)

func TestNamespace_New(t *testing.T) {
	gt := NewGomegaWithT(t)

	db, cleanup := newKVDB(t)
	defer cleanup()

	logger := zap.NewExample()

	ns := New(logger, db)
	gt.Expect(ns.Logger).To(Equal(logger))
	gt.Expect(ns.TxRepo).NotTo(BeNil())
	gt.Expect(ns.SubmitService).NotTo(BeNil())
}

func newKVDB(t *testing.T) (*store.LevelDBKV, func()) {
	path, cleanup := tested.TempDir(t, "", "level")

	db, err := store.NewLevelDB(path)
	NewGomegaWithT(t).Expect(err).NotTo(HaveOccurred())

	return db, cleanup
}
