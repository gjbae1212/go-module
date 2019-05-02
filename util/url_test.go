package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	assert := assert.New(t)

	url1 := "http://naver.com/aa/bb?cc=dd&ee=ff#fragment"
	schema, domain, port, path, query, fragment, err := ParseURL(url1)
	assert.NoError(err)
	assert.Equal("http", schema)
	assert.Equal("naver.com", domain)
	assert.Equal("80", port)
	assert.Equal("/aa/bb", path)
	assert.Equal("cc=dd&ee=ff", query)
	assert.Equal("fragment", fragment)

	url2 := "cc.com:8080/aa/bb"
	schema, domain, port, path, query, fragment, err = ParseURL(url2)
	assert.NoError(err)
	assert.Equal("http", schema)
	assert.Equal("cc.com", domain)
	assert.Equal("8080", port)
	assert.Equal("/aa/bb", path)

	url3:= "https://naver.com/aa/bb?cc=dd&ee=ff#fragment"
	schema, domain, port, path, query, fragment, err = ParseURL(url3)
	assert.Equal("https", schema)
	assert.Equal("naver.com", domain)
	assert.Equal("443", port)
	assert.Equal("/aa/bb", path)
	assert.Equal("cc=dd&ee=ff", query)
	assert.Equal("fragment", fragment)
}
