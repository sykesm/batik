// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpcapi

import (
	"context"
	"crypto"
	"testing"

	. "github.com/onsi/gomega"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestEncode(t *testing.T) {
	gt := NewGomegaWithT(t)

	testTx := newTestTransaction()
	itx, err := transaction.New(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())
	req := &txv1.EncodeRequest{Transaction: testTx}

	encodeSvc := &EncodeService{}
	resp, err := encodeSvc.Encode(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp.Txid).To(Equal(itx.ID.Bytes()))

	expectedEncoded, err := protomsg.MarshalDeterministic(testTx)
	gt.Expect(resp.EncodedTransaction).To(Equal(expectedEncoded))
}
