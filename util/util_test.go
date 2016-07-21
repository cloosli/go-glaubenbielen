package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeText(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(NormalizeText("Bürgisweiher Sprint 125 éàè!!"), "BurgisweiherSprint125eae")
	assert.Equal(NormalizeText("123éü+*ç"), "123euc")
	assert.Equal(NormalizeText("file.csv"), "filecsv")
}

func TestFloatToString(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("1.00000000", FloatToString(1.0000000000000))
	assert.Equal("-0.00000000", FloatToString(-0.00000000000000001))
	assert.Equal("47.12345679", FloatToString(47.123456789123456))
}
