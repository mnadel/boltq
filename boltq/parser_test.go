package boltq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanParseSimpleSelect(t *testing.T) {
	p := NewParser(strings.NewReader("select a from b"))

	parsed, err := p.ParseSelect()
	assert.NoError(t, err, "error parsing")

	assert.Equal(t, 1, len(parsed.Fields))
	field := parsed.Fields[0]
	assert.Equal(t, "a", field)

	assert.Equal(t, 1, len(parsed.BucketPath))
	bucket := parsed.BucketPath[0]
	assert.Equal(t, "b", bucket)
}

func TestCanParseMultipleFields(t *testing.T) {
	p := NewParser(strings.NewReader("select a,b from c"))

	parsed, err := p.ParseSelect()
	assert.NoError(t, err, "error parsing")

	assert.Equal(t, 2, len(parsed.Fields))
	field := parsed.Fields[0]
	assert.Equal(t, "a", field)
	field = parsed.Fields[1]
	assert.Equal(t, "b", field)
}

func TestCanParseMultipleBuckets(t *testing.T) {
	p := NewParser(strings.NewReader("select a,b from b1/b2/c"))

	parsed, err := p.ParseSelect()
	assert.NoError(t, err, "error parsing")

	assert.Equal(t, 3, len(parsed.BucketPath))
	bucket := parsed.BucketPath[0]
	assert.Equal(t, "b1", bucket)
	bucket = parsed.BucketPath[1]
	assert.Equal(t, "b2", bucket)
	bucket = parsed.BucketPath[2]
	assert.Equal(t, "c", bucket)
}
