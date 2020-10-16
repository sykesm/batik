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
	gt.Expect(resp.Txid).To(Equal(fromHex(t, "77dc6e1729583cf7f1db9863b34a8951a3bb9369ab4cf0a86340ea92a8514cf5")))

	expectedEncoded, err := protomsg.MarshalDeterministic(testTx)
	gt.Expect(resp.EncodedTransaction).To(Equal(expectedEncoded))
}
