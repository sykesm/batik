// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"go.uber.org/zap"
)

// syncLogger returns an exit handler function for syncing a *zap.Logger.
func syncLogger(logger *zap.Logger) func() {
	return func() {
		logger.Sync()
	}
}
