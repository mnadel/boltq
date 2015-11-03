package boltq

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanLexSimpleStringUpdate(t *testing.T) {
	s := NewScanner(strings.NewReader("update bucket set key1 = 'new_value'"))

	expectedTokens := []Token{UPDATE, WS, IDENT, WS, SET, WS, IDENT, WS, EQUALS, WS, QUOTE, IDENT, QUOTE, EOF}

	for i := range expectedTokens {
		tok, _ := s.Scan()
		assert.Equal(t, expectedTokens[i], tok, fmt.Sprintf("i = %d", i))
	}
}

func TestCanLexSimpleNumericUpdate(t *testing.T) {
	s := NewScanner(strings.NewReader("update bucket set key1 = 2.2"))

	expectedTokens := []Token{UPDATE, WS, IDENT, WS, SET, WS, IDENT, WS, EQUALS, WS, IDENT, EOF}

	for i := range expectedTokens {
		tok, _ := s.Scan()
		assert.Equal(t, expectedTokens[i], tok, fmt.Sprintf("i = %d", i))
	}
}

func TestCanLexNestedBucketUpdate(t *testing.T) {
	s := NewScanner(strings.NewReader("update bucket/subbucket set key1 = 'new_value'"))

	expectedTokens := []Token{UPDATE, WS, IDENT, BUCKETSEP, IDENT, WS, SET, WS, IDENT, WS, EQUALS, WS, QUOTE, IDENT, QUOTE, EOF}

	for i := range expectedTokens {
		tok, _ := s.Scan()
		assert.Equal(t, expectedTokens[i], tok, fmt.Sprintf("i = %d", i))
	}
}

func TestCanLexMultiValueUpdate(t *testing.T) {
	s := NewScanner(strings.NewReader("update bucket set key1 = 'new_value', key2 = 3.14"))

	expectedTokens := []Token{UPDATE, WS, IDENT, WS, SET, WS, IDENT, WS, EQUALS, WS, QUOTE, IDENT, QUOTE, COMMA, WS, IDENT, WS, EQUALS, WS, IDENT, EOF}

	for i := range expectedTokens {
		tok, _ := s.Scan()
		assert.Equal(t, expectedTokens[i], tok, fmt.Sprintf("i = %d", i))
	}
}
