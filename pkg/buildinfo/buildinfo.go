// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package buildinfo contains build information that is provided during
// the production build process.
package buildinfo

import "time"

var (
	Version   = "dev"     // Version is the version of the program that was built.
	GitCommit = "unknown" // GitCommit is the git commit that was used when the program was built.
	Timestamp = "unknown" // Timestamp is the time when the program build process was started.
)

// FullVersion returns the concatenation of the version and the git commit
// hash.
func FullVersion() string {
	if len(GitCommit) > 0 {
		return Version + "-" + GitCommit
	}
	return Version
}

// Built returns the build Timestamp as a time.Time. If parsing the time fails,
// a zero time is returned.
func Built() time.Time {
	// ignore parse failure and return a zero time
	built, _ := time.Parse(time.RFC3339, Timestamp)
	return built
}
