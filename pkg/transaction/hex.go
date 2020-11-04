// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import "encoding/hex"

func toHexString(b []byte) string {
	if len(b) <= 32 {
		return hex.EncodeToString(b)
	}
	return hex.EncodeToString(b[0:15]) + "...." + hex.EncodeToString(b[len(b)-15:])
}
