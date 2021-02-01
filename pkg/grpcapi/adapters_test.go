// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpcapi

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sykesm/batik/pkg/namespace"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/submit"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestAdapters_MapAdapter(t *testing.T) {
	gt := NewGomegaWithT(t)

	repoPtr := &store.TransactionRepository{}
	submitPtr := &submit.Service{}

	adapter := NamespaceMapAdapter(map[string]*namespace.Namespace{
		"present": {
			TxRepo:        repoPtr,
			SubmitService: submitPtr,
		},
	})

	store := adapter.Repository("present")
	gt.Expect(store).To(Equal(repoPtr))

	submit := adapter.Submitter("present")
	gt.Expect(submit).To(Equal(submitPtr))

	missingStore := adapter.Repository("missing")
	gt.Expect(missingStore).To(Equal(notFoundRepository("missing")))

	missingSubmit := adapter.Submitter("missing")
	gt.Expect(missingSubmit).To(Equal(notFoundSubmitter("missing")))
}

func TestAdapters_NotFoundSubmitter(t *testing.T) {
	gt := NewGomegaWithT(t)

	nfs := notFoundSubmitter("missing")
	err := nfs.Submit(nil, nil)
	gt.Expect(err).To(MatchError("bad namespace \"missing\": namespace not found"))
}

func TestAdapters_NotFoundRepository(t *testing.T) {
	gt := NewGomegaWithT(t)

	nfr := notFoundRepository("missing")

	err := nfr.PutTransaction(nil)
	gt.Expect(err).To(MatchError("bad namespace \"missing\": namespace not found"))

	_, err = nfr.GetTransaction(nil)
	gt.Expect(err).To(MatchError("bad namespace \"missing\": namespace not found"))

	err = nfr.PutState(nil)
	gt.Expect(err).To(MatchError("bad namespace \"missing\": namespace not found"))

	_, err = nfr.GetState(transaction.StateID{}, false)
	gt.Expect(err).To(MatchError("bad namespace \"missing\": namespace not found"))
}
