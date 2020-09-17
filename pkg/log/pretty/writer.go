// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"io"
)

type Writer struct {
	dst io.Writer
}

// NewWriter returns a pretty Writer that prettifys logfmt lines.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		dst: w,
	}
}

// Write prettifys logfmt lines and writes them onto the writer's dst.
// If the lines aren't logfmt, it will simply write them out with no
// prettification.
func (w *Writer) Write(p []byte) (n int, err error) {
	if logfmtEntry := UnmarshalLogfmt(p); logfmtEntry != nil {
		return w.dst.Write(logfmtEntry.Prettify())
	}

	return w.dst.Write(p)
}
