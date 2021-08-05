package binder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse_short(t *testing.T) {
	v := ParseTag("json")
	r := map[string]string{
		"json": "",
	}
	assert.Equal(t, r, v)
}

func TestParse_short_long(t *testing.T) {
	// 'json' and 'arg' will have the same value
	v := ParseTag("json yaml arg:\"abc\"")
	r := map[string]string{
		"json": "abc",
		"yaml": "abc",
		"arg":  "abc",
	}
	assert.Equal(t, r, v)
}

func TestParse_long_short(t *testing.T) {
	v := ParseTag("arg:\"abc\" json yaml")
	r := map[string]string{
		"json": "",
		"yaml": "",
		"arg":  "abc",
	}
	assert.Equal(t, r, v)
}

func TestParse_quotedLong(t *testing.T) {
	r := map[string]string{
		"json": "",
		"yaml": "",
		"arg":  "abc",
	}
	v := ParseTag("arg:abc json yaml")
	assert.Equal(t, r, v)
	v = ParseTag("arg:\"abc\" json yaml")
	assert.Equal(t, r, v)
}

func TestParse_panic_1(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	ParseTag("xml arg:`abc` json yaml")
}

func TestParse_panic_2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	ParseTag("arg:'abc' json")
}
