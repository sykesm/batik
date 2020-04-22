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

type argScanner struct {
	state      state
	text       string
	pos        int
	buf        *bytes.Buffer
	incomplete bool
	args       []string
}

func (s *argScanner) ScanLine(text string) (bool, error) {
	s.text = text
	s.pos = 0
	s.incomplete = false

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
				s.incomplete = true
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

func (s *argScanner) Args() ([]string, error) {
	if s.state == doubleQuoteState {
		return nil, errors.New("scanner: double quoted string not terminated")
	}
	if s.state == singleQuoteState {
		return nil, errors.New("scanner: single quoted string not terminated")
	}
	if s.incomplete {
		return nil, errors.New("scanner: incomplete scan")
	}
	return s.args, nil
}

func (s *argScanner) Reset() {
	s.state = initialState
	s.text = ""
	s.pos = 0
	s.buf = nil
	s.incomplete = false
	s.args = nil
}

func (s *argScanner) completeToken() {
	if s.buf == nil {
		return
	}

	s.args = append(s.args, s.buf.String())
	s.buf = nil
}

func (s *argScanner) writeRune(ch rune) error {
	if s.buf == nil {
		s.buf = bytes.NewBuffer(nil)
	}

	_, err := s.buf.WriteRune(ch)
	return err
}

func (s *argScanner) next() rune {
	if s.pos >= len(s.text) {
		return eol
	}

	ch, l := utf8.DecodeRuneInString(s.text[s.pos:])
	s.pos += l

	return ch // may be RuneError
}
