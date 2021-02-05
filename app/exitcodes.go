// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

const (
	exitOkay = iota
	exitShellSetupFailed
	exitCommandNotFound
	exitConfigLoadFailed
	exitLoggerCreateFailed
	exitErrLoggerCreateFailed
	exitServerCreateFailed
	exitServerStartFailed
	exitServerStatusFailed
	exitAppShutdownFailed
	exitChangeLogspecFailed
	exitConfigEncodeFailed
)
