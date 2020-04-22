// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package repl

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		args  []string
	}{
		{
			name:  "SingleWord",
			input: []string{"word"},
			args:  []string{"word"},
		},
		{
			name:  "MultipleWords",
			input: []string{"word1 word2  word3"},
			args:  []string{"word1", "word2", "word3"},
		},
		{
			name:  "TrailingSpace",
			input: []string{"word1 "},
			args:  []string{"word1"},
		},
		{
			name:  "MultiLineWords",
			input: []string{"word1\\", "word2 \\", "word3 "},
			args:  []string{"word1word2", "word3"},
		},
		{
			name:  "MultiWordDoubleQuotes",
			input: []string{`"multi word string"`},
			args:  []string{`multi word string`},
		},
		{
			name:  "MultiLineDoubleQuotes",
			input: []string{`"line one`, `line two `, `line three"`},
			args:  []string{"line one\nline two \nline three"},
		},
		{
			name:  "MultiWordSingleQuotes",
			input: []string{`'multi word string'`},
			args:  []string{`multi word string`},
		},
		{
			name:  "MultiLineSingleQuotes",
			input: []string{`'line one`, `line two `, `line three'`},
			args:  []string{"line one\nline two \nline three"},
		},
		{
			name:  "ConcatDoubleQuotes",
			input: []string{`"foo""bar"`},
			args:  []string{"foobar"},
		},
		{
			name:  "ConcatSingleQuotes",
			input: []string{`'foo''bar'`},
			args:  []string{"foobar"},
		},
		{
			name:  "ConcatMixedQuotes",
			input: []string{`'foo'"bar"'baz'boo`},
			args:  []string{"foobarbazboo"},
		},
		{
			name:  "ConcatLinesMixedQuotes",
			input: []string{`"foo"'bar`, `"baz'quo`},
			args:  []string{"foobar\nbazquo"},
		},
		{
			name:  "EscapedSingleQuote",
			input: []string{`That\'s the spirit\!`},
			args:  []string{"That's", "the", "spirit!"},
		},
		{
			name:  "EscapeInSingleQuotes",
			input: []string{`'\n'`},
			args:  []string{`\n`},
		},
		{
			name:  "EscapeInDoubleQuotes",
			input: []string{`"\\Paid: \$12.00 \k"`},
			args:  []string{`\Paid: $12.00 \k`},
		},
		{
			name:  "EmptyQuotedString",
			input: []string{`one '' "" four`},
			args:  []string{"one", "", "", "four"},
		},
		{
			name:  "EscapedSpace",
			input: []string{`foo \ \`, `bar`},
			args:  []string{"foo", " bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			scanner := &argScanner{}
			for i, input := range tt.input {
				more, err := scanner.ScanLine(input)
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(more).To(Equal(i != len(tt.input)-1))
			}

			args, err := scanner.Args()
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(args).To(Equal(tt.args))
		})
	}
}

func TestScannerErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{
			name:  "IncompleteDoubleQuoted",
			input: `"I never finish anyth`,
			err:   "scanner: double quoted string not terminated",
		},
		{
			name:  "IncompleteSingleQuoted",
			input: `'I never "finish" anyth`,
			err:   "scanner: single quoted string not terminated",
		},
		{
			name:  "EndWithLineContinuation",
			input: `line one\`,
			err:   "scanner: incomplete scan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			scanner := &argScanner{}
			_, err := scanner.ScanLine(tt.input)
			gt.Expect(err).NotTo(HaveOccurred())

			_, err = scanner.Args()
			gt.Expect(err).To(MatchError(tt.err))
		})
	}
}

func TestScannerScanLineMultipleCalls(t *testing.T) {
	gt := NewGomegaWithT(t)

	scanner := &argScanner{}
	more, err := scanner.ScanLine("word1")
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(more).To(BeFalse())

	args, err := scanner.Args()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(args).To(Equal([]string{"word1"}))

	more, err = scanner.ScanLine("word2")
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(more).To(BeFalse())

	args, err = scanner.Args()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(args).To(Equal([]string{"word1", "word2"}))
}

func TestScannerArgsMultipleCalls(t *testing.T) {
	gt := NewGomegaWithT(t)

	scanner := &argScanner{}
	more, err := scanner.ScanLine("these are words")
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(more).To(BeFalse())

	tok1, err := scanner.Args()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(tok1).To(Equal([]string{"these", "are", "words"}))

	tok2, err := scanner.Args()
	gt.Expect(tok2).To(Equal(tok1))
}

func TestScannerReset(t *testing.T) {
	gt := NewGomegaWithT(t)

	scanner := &argScanner{}
	more, err := scanner.ScanLine("these are words")
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(more).To(BeFalse())

	args, err := scanner.Args()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(args).To(Equal([]string{"these", "are", "words"}))

	scanner.Reset()
	args, err = scanner.Args()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(args).To(BeEmpty())

	more, err = scanner.ScanLine("these are words")
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(more).To(BeFalse())

	args, err = scanner.Args()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(args).To(Equal([]string{"these", "are", "words"}))
}
