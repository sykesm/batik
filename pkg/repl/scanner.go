// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package repl

import (
	"bytes"
	"errors"
	"unicode/utf8"
)

const eol = -1

type state uint8

const (
	initialState state = iota
	doubleQuoteState
	singleQuoteState
)

// TODO(mjs): unexport
type Scanner struct {
	state        state
	text         string
	pos          int
	buf          *bytes.Buffer
	moreExpected bool
	tokens       []string
}

func (s *Scanner) ScanLine(text string) (bool, error) {
	s.text = text
	s.pos = 0
	s.moreExpected = false

	for ch := s.next(); ch != eol; ch = s.next() {
		if ch == utf8.RuneError {
			return false, errors.New("scanner: rune error")
		}

		switch ch {
		case '\t', ' ':
			if s.state == initialState {
				s.completeToken()
			} else {
				s.writeRune(ch)
			}

		case '"':
			if s.state == initialState {
				s.state = doubleQuoteState
				if s.buf == nil {
					s.buf = bytes.NewBuffer(nil)
				}
			} else if s.state == doubleQuoteState {
				s.state = initialState
			}

		case '\'':
			if s.state == initialState {
				s.state = singleQuoteState
				if s.buf == nil {
					s.buf = bytes.NewBuffer(nil)
				}
			} else if s.state == singleQuoteState {
				s.state = initialState
			}

		case '\\':
			ch = s.next()
			if ch == eol { // line continuation
				s.moreExpected = true
				return true, nil
			}

			if s.state == doubleQuoteState {
				switch ch {
				case '$', '`', '"', '\\':
					// The backslash retains its special meaning only when followed by
					// one of the following characters: ‘$’, ‘`’, ‘"’, ‘\’, or newline.
					// https://www.gnu.org/software/bash/manual/html_node/Double-Quotes.html
				default:
					if err := s.writeRune('\\'); err != nil {
						return false, err
					}
				}
			}
			if s.state == singleQuoteState {
				if err := s.writeRune('\\'); err != nil {
					return false, err
				}
			}
			fallthrough

		default:
			if err := s.writeRune(ch); err != nil {
				return false, err
			}
		}
	}

	if s.state != initialState {
		s.writeRune('\n')
	}
	if s.state == initialState {
		s.completeToken()
	}

	return s.state != initialState, nil
}

func (s *Scanner) completeToken() {
	if s.buf == nil {
		return
	}

	s.tokens = append(s.tokens, s.buf.String())
	s.buf = nil
}

func (s *Scanner) writeRune(ch rune) error {
	if s.buf == nil {
		s.buf = bytes.NewBuffer(nil)
	}

	_, err := s.buf.WriteRune(ch)
	return err
}

func (s *Scanner) next() rune {
	if s.pos >= len(s.text) {
		return eol
	}

	ch, l := utf8.DecodeRuneInString(s.text[s.pos:])
	s.pos += l

	return ch // may be RuneError
}

func (s *Scanner) Tokens() ([]string, error) {
	if s.state == doubleQuoteState {
		return nil, errors.New("scanner: double quoted string not terminated")
	}
	if s.state == singleQuoteState {
		return nil, errors.New("scanner: single quoted string not terminated")
	}
	if s.moreExpected {
		return nil, errors.New("scanner: incomplete scan")
	}
	return s.tokens, nil
}

func (s *Scanner) Reset() {
	s.state = initialState
	s.text = ""
	s.pos = 0
	s.buf = nil
	s.moreExpected = false
	s.tokens = nil
}
