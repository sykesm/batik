package transaction

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"google.golang.org/protobuf/proto"
)

func TestEncodeService(t *testing.T) {
	gt := NewGomegaWithT(t)

	testTx := newTestTransaction()

	req := &tb.EncodedTransactionRequest{
		Transaction: testTx,
	}

	encodeSvc := &EncodeService{}
	response, err := encodeSvc.EncodedTransaction(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(response.Txid).To(Equal(fromHex(t, "53e33ae87fb6cf2e4aaaabcdae3a93d578d9b7366e905dfff0446356774f726f")))

	expectedEncoded, err := proto.MarshalOptions{Deterministic: true}.Marshal(testTx)
	gt.Expect(response.EncodedTransaction).To(Equal(expectedEncoded))
}
