package boltq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanUpdateInteger(t *testing.T) {
	p := NewParser(strings.NewReader("update a set b = 3"))

	parsed, err := p.ParseUpdate()
	assert.NoError(t, err, "error parsing")

	assert.Equal(t, 1, len(parsed.BucketPath))
	bucket := parsed.BucketPath[0]
	assert.Equal(t, "a", bucket)

	assert.Equal(t, 1, len(parsed.Fields))
	assert.EqualValues(t, 3, parsed.Fields["b"])
}

func TestCanUpdateFloat(t *testing.T) {
	p := NewParser(strings.NewReader("update a set b = 3.14"))

	parsed, err := p.ParseUpdate()
	assert.NoError(t, err, "error parsing")

	assert.Equal(t, 1, len(parsed.BucketPath))
	bucket := parsed.BucketPath[0]
	assert.Equal(t, "a", bucket)

	assert.Equal(t, 1, len(parsed.Fields))
	assert.EqualValues(t, 3.14, parsed.Fields["b"])
}

func TestCanUpdateString(t *testing.T) {
	p := NewParser(strings.NewReader("update a set b = 'new_value'"))

	parsed, err := p.ParseUpdate()
	assert.NoError(t, err, "error parsing")

	assert.Equal(t, 1, len(parsed.BucketPath))
	bucket := parsed.BucketPath[0]
	assert.Equal(t, "a", bucket)

	assert.Equal(t, 1, len(parsed.Fields))
	assert.Equal(t, "new_value", parsed.Fields["b"])
}

func TestCanUpdateBucketPath(t *testing.T) {
	p := NewParser(strings.NewReader("update a/b/c set d = e"))

	parsed, err := p.ParseUpdate()
	assert.NoError(t, err, "error parsing")

	assert.Equal(t, 3, len(parsed.BucketPath))
	assert.Equal(t, "a", parsed.BucketPath[0])
	assert.Equal(t, "b", parsed.BucketPath[1])
	assert.Equal(t, "c", parsed.BucketPath[2])

	assert.Equal(t, 1, len(parsed.Fields))
	assert.Equal(t, "e", parsed.Fields["d"])
}

func TestCanUpdateMultipleKeys(t *testing.T) {
	p := NewParser(strings.NewReader("update a set b = 3.14, c = dog, d = 'foobar'"))

	parsed, err := p.ParseUpdate()
	assert.NoError(t, err, "error parsing")

	assert.Equal(t, 1, len(parsed.BucketPath))
	assert.Equal(t, "a", parsed.BucketPath[0])

	assert.Equal(t, 3, len(parsed.Fields))
	assert.EqualValues(t, 3.14, parsed.Fields["b"])
	assert.Equal(t, "dog", parsed.Fields["c"])
	assert.Equal(t, "foobar", parsed.Fields["d"])
}
