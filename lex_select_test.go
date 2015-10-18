package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanLexSimpleSelect(t *testing.T) {
	s := NewScanner(strings.NewReader("select * from bucket"))

	expectedTokens := []Token{SELECT, WS, ASTERISK, WS, FROM, WS, IDENT, EOF}

	for i := range expectedTokens {
		tok, lit := s.Scan()

		switch tok {
		case IDENT:
			assert.Equal(t, expectedTokens[i], tok)
			assert.Equal(t, "bucket", lit)
		default:
			assert.Equal(t, expectedTokens[i], tok)
		}
	}
}

func TestCanLexNestedBucketSelect(t *testing.T) {
	s := NewScanner(strings.NewReader("select * from bucket/subbucket"))

	expectedTokens := []Token{SELECT, WS, ASTERISK, WS, FROM, WS, IDENT, BUCKETSEP, IDENT, EOF}

	for i := range expectedTokens {
		tok, lit := s.Scan()

		switch tok {
		case IDENT:
			assert.Equal(t, expectedTokens[i], tok)
			if i == 6 {
				assert.Equal(t, "bucket", lit)
			} else if i == 8 {
				assert.Equal(t, "subbucket", lit)
			}
		default:
			assert.Equal(t, expectedTokens[i], tok)
		}
	}
}

func TestCanLexMultiSelect(t *testing.T) {
	s := NewScanner(strings.NewReader("select a,b from bucket"))

	s.Scan() // select
	s.Scan() // ws

	tok, lit := s.Scan()
	assert.Equal(t, IDENT, tok)
	assert.Equal(t, "a", lit)

	tok, _ = s.Scan()
	assert.Equal(t, COMMA, tok)

	tok, lit = s.Scan()
	assert.Equal(t, IDENT, tok)
	assert.Equal(t, "b", lit)
}
