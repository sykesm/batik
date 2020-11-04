// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"testing"

	. "github.com/onsi/gomega"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	. "github.com/sykesm/batik/pkg/tested/matcher"
)

func TestStateConversion(t *testing.T) {
	stateID := StateID{
		TxID:        ID([]byte("transaction-id-0")),
		OutputIndex: 1,
	}
	protoState := &txv1.State{
		Info: &txv1.StateInfo{
			Kind: "state-kind-0",
			Owners: []*txv1.Party{
				{PublicKey: []byte("owner-1")},
				{PublicKey: []byte("owner-2")},
			},
		},
		State: []byte("state-data"),
	}
	state := &State{
		ID: stateID,
		StateInfo: &StateInfo{
			Kind: "state-kind-0",
			Owners: []*Party{
				{PublicKey: []byte("owner-1")},
				{PublicKey: []byte("owner-2")},
			},
		},
		Data: []byte("state-data"),
	}

	t.Run("ToState", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToState(nil, stateID.TxID, stateID.OutputIndex)).To(BeNil())
		gt.Expect(ToState(protoState, stateID.TxID, stateID.OutputIndex)).To(Equal(state))
	})

	t.Run("ToStates", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToStates(stateID.TxID)).To(BeEmpty())
		gt.Expect(ToStates(stateID.TxID, nil)).To(Equal([]*State{nil}))

		state0, state1 := *state, *state
		state0.ID.OutputIndex, state1.ID.OutputIndex = 0, 1
		gt.Expect(ToStates(stateID.TxID, protoState, protoState)).To(Equal([]*State{&state0, &state1}))
	})

	t.Run("FromState", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FromState(nil)).To(BeNil())
		gt.Expect(FromState(state)).To(ProtoEqual(protoState))
	})
}

func TestStateInfoConversion(t *testing.T) {
	protoStateInfo := &txv1.StateInfo{
		Kind: "state-kind-0",
		Owners: []*txv1.Party{
			{PublicKey: []byte("owner-1")},
			{PublicKey: []byte("owner-2")},
		},
	}
	stateInfo := &StateInfo{
		Kind: "state-kind-0",
		Owners: []*Party{
			{PublicKey: []byte("owner-1")},
			{PublicKey: []byte("owner-2")},
		},
	}

	t.Run("ToStateInfo", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToStateInfo(nil)).To(BeNil())
		gt.Expect(ToStateInfo(protoStateInfo)).To(Equal(stateInfo))
	})

	t.Run("FromStateInfo", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FromStateInfo(nil)).To(BeNil())
		gt.Expect(FromStateInfo(stateInfo)).To(ProtoEqual(protoStateInfo))
	})
}

func TestStateIDConversion(t *testing.T) {
	txid := NewID([]byte("transaction-id-0"))
	protoStateRef := &txv1.StateReference{
		Txid:        txid.Bytes(),
		OutputIndex: 11,
	}
	stateID := &StateID{
		TxID:        txid,
		OutputIndex: 11,
	}

	t.Run("ToStateID", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToStateID(nil)).To(BeNil())
		gt.Expect(ToStateID(protoStateRef)).To(Equal(stateID))
	})

	t.Run("ToStateIDs", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToStateIDs()).To(BeEmpty())
		gt.Expect(ToStateIDs(nil)).To(Equal([]*StateID{nil}))
		gt.Expect(ToStateIDs(protoStateRef, protoStateRef)).To(Equal([]*StateID{stateID, stateID}))
	})

	t.Run("FromStateID", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FromStateID(nil)).To(BeNil())
		gt.Expect(FromStateID(stateID)).To(ProtoEqual(protoStateRef))
	})

	t.Run("FromStateIDs", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FromStateIDs()).To(BeEmpty())
		gt.Expect(FromStateIDs(nil)).To(Equal([]*txv1.StateReference{nil}))
		gt.Expect(FromStateIDs(stateID, stateID)).To(ConsistOf(
			ProtoEqual(protoStateRef),
			ProtoEqual(protoStateRef),
		))
	})
}

func TestPartyConversion(t *testing.T) {
	protoParties := []*txv1.Party{
		{PublicKey: []byte("public-key-0")},
		{PublicKey: []byte("public-key-1")},
	}
	parties := []*Party{
		{PublicKey: []byte("public-key-0")},
		{PublicKey: []byte("public-key-1")},
	}

	t.Run("ToParty", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToParty(nil)).To(BeNil())
		gt.Expect(ToParty(protoParties[0])).To(Equal(parties[0]))
	})

	t.Run("ToParties", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToParties()).To(BeEmpty())
		gt.Expect(ToParties(nil)).To(Equal([]*Party{nil}))
		gt.Expect(ToParties(protoParties...)).To(Equal(parties))
	})

	t.Run("FromParty", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FromParty(nil)).To(BeNil())
		gt.Expect(FromParty(parties[0])).To(ProtoEqual(protoParties[0]))
	})

	t.Run("FromParties", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FromParties()).To(BeEmpty())
		gt.Expect(FromParties(nil)).To(Equal([]*txv1.Party{nil}))
		gt.Expect(FromParties(parties...)).To(ConsistOf(
			ProtoEqual(protoParties[0]),
			ProtoEqual(protoParties[1]),
		))
	})
}

func TestParameterConversion(t *testing.T) {
	protoParameters := []*txv1.Parameter{
		{Name: "key-1", Value: []byte("value-1")},
		{Name: "key-2", Value: []byte("value-2")},
	}
	parameters := []*Parameter{
		{Name: "key-1", Value: []byte("value-1")},
		{Name: "key-2", Value: []byte("value-2")},
	}

	t.Run("ToParameter", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToParameter(nil)).To(BeNil())
		gt.Expect(ToParameter(protoParameters[0])).To(Equal(parameters[0]))
	})

	t.Run("ToParameters", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToParameters()).To(BeEmpty())
		gt.Expect(ToParameters(nil)).To(Equal([]*Parameter{nil}))
		gt.Expect(ToParameters(protoParameters...)).To(Equal(parameters))
	})

	t.Run("FromParameter", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FromParameter(nil)).To(BeNil())
		gt.Expect(FromParameter(parameters[0])).To(ProtoEqual(protoParameters[0]))
	})

	t.Run("FromParameters", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FromParameters()).To(BeEmpty())
		gt.Expect(FromParameters(nil)).To(Equal([]*txv1.Parameter{nil}))
		gt.Expect(FromParameters(parameters...)).To(ConsistOf(
			ProtoEqual(protoParameters[0]),
			ProtoEqual(protoParameters[1]),
		))
	})
}

func TestSignatureConversion(t *testing.T) {
	protoSignatures := []*txv1.Signature{
		{PublicKey: []byte("public-key-1"), Signature: []byte("signature-1")},
		{PublicKey: []byte("public-key-2"), Signature: []byte("signature-2")},
	}
	signatures := []*Signature{
		{PublicKey: []byte("public-key-1"), Signature: []byte("signature-1")},
		{PublicKey: []byte("public-key-2"), Signature: []byte("signature-2")},
	}

	t.Run("ToSignature", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToSignature(nil)).To(BeNil())
		gt.Expect(ToSignature(protoSignatures[0])).To(Equal(signatures[0]))
	})

	t.Run("ToSignatures", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(ToSignatures()).To(BeEmpty())
		gt.Expect(ToSignatures(nil)).To(Equal([]*Signature{nil}))
		gt.Expect(ToSignatures(protoSignatures...)).To(Equal(signatures))
	})

	t.Run("FromSignatures", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(FromSignatures()).To(BeEmpty())
		gt.Expect(FromSignatures(nil)).To(Equal([]*txv1.Signature{nil}))
		gt.Expect(FromSignatures(signatures...)).To(ConsistOf(
			ProtoEqual(protoSignatures[0]),
			ProtoEqual(protoSignatures[1]),
		))
	})
}
