// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package totalorder

import (
	"github.com/sykesm/batik/pkg/transaction"
)

// TXIDAndHMAC is a pair of txid, and an HMAC computed using
// the namespace name, a secret, and the txid.  This allows for
// namespace members to detect transactions for their namespace
// while other namespace members can only discern it is not
// for a namespcae they care about.
type TXIDAndHMAC struct {
	ID   transaction.ID
	HMAC []byte
}

func (t *TXIDAndHMAC) serialize() []byte {
	if len(t.ID) != 32 || len(t.HMAC) != 32 {
		// XXX we should probably define a better serialization?
		panic("we are serializing based on assumed offsets")
	}

	return append(append([]byte{}, t.ID...), t.HMAC...)
}

func txidAndHMACFromBytes(b []byte) TXIDAndHMAC {
	if len(b) != 64 {
		panic("unexpected size for serialization, should be fixed 64 bytes")
	}

	return TXIDAndHMAC{
		ID:   b[:32],
		HMAC: b[32:],
	}
}
