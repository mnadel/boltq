// adapted from: http://blog.gopheracademy.com/advent-2014/parsers-lexers/

package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const (
	ILLEGAL Token = iota
	EOF
	WS

	IDENT

	ASTERISK
	BUCKETSEP
	COMMA
	QUOTE
	EQUALS

	SELECT
	FROM
	UPDATE
	DELETE
	SET
	WHERE
)

type Token int

type Scanner struct {
	r *bufio.Reader
}

var eof = rune(0)

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}

	return ch
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isIdentChar(ch) {
		s.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case '/':
		return BUCKETSEP, string(ch)
	case ',':
		return COMMA, string(ch)
	case '\'':
		return QUOTE, string(ch)
	case '=':
		return EQUALS, string(ch)
	}

	return ILLEGAL, string(ch)
}

func isIdentChar(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '.' || ch == ':'
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isIdentChar(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	switch strings.ToUpper(buf.String()) {
	case "SELECT":
		return SELECT, buf.String()
	case "FROM":
		return FROM, buf.String()
	case "UPDATE":
		return UPDATE, buf.String()
	case "DELETE":
		return DELETE, buf.String()
	case "WHERE":
		return WHERE, buf.String()
	case "SET":
		return SET, buf.String()
	}

	// Otherwise return as a regular identifier.
	return IDENT, buf.String()
}
