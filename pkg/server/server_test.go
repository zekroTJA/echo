package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseContentType(t *testing.T) {
	var typ, rest string

	typ, rest = parseContentType("application/json")
	assert.Equal(t, "application/json", typ)
	assert.Equal(t, "", rest)

	typ, rest = parseContentType("text/html; charset=utf-8")
	assert.Equal(t, "text/html", typ)
	assert.Equal(t, "charset=utf-8", rest)

	typ, rest = parseContentType("TEXT/HTML;charset=utf-8")
	assert.Equal(t, "text/html", typ)
	assert.Equal(t, "charset=utf-8", rest)

	typ, rest = parseContentType("")
	assert.Equal(t, "", typ)
	assert.Equal(t, "", rest)

	typ, rest = parseContentType("foo;")
	assert.Equal(t, "foo", typ)
	assert.Equal(t, "", rest)

	typ, rest = parseContentType("; bar baz")
	assert.Equal(t, "", typ)
	assert.Equal(t, "bar baz", rest)

	typ, rest = parseContentType(";")
	assert.Equal(t, "", typ)
	assert.Equal(t, "", rest)
}
