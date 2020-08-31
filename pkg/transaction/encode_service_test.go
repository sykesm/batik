// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/protomsg"
)

func TestEncodeService(t *testing.T) {
	gt := NewGomegaWithT(t)

	testTx := newTestTransaction()
	req := &tb.EncodeTransactionRequest{Transaction: testTx}

	encodeSvc := &EncodeService{}
	resp, err := encodeSvc.EncodeTransaction(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp.Txid).To(Equal(fromHex(t, "53e33ae87fb6cf2e4aaaabcdae3a93d578d9b7366e905dfff0446356774f726f")))

	expectedEncoded, err := protomsg.MarshalDeterministic(testTx)
	gt.Expect(resp.EncodedTransaction).To(Equal(expectedEncoded))
}
