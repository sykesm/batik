// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"fmt"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
)

type StateID struct {
	TxID        ID
	OutputIndex uint64
}

func (sid StateID) String() string {
	return fmt.Sprintf("%s:%016x", sid.TxID, sid.OutputIndex)
}

func (sid StateID) Equals(that StateID) bool {
	if sid.OutputIndex == that.OutputIndex && sid.TxID.Equals(that.TxID) {
		return true
	}
	return false
}

type State struct {
	ID        StateID
	StateInfo *txv1.StateInfo
	Data      []byte
}
