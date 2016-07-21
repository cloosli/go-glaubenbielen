package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNormalizeText(t *testing.T) {
	assert := assert.New(t)
	name := "Bürgisweiher Sprint 125 éàè!!"
	fmt.Println(strings.Title(name))
	fmt.Println(strings.ToTitle(name))
	fmt.Println(NormalizeText(name))
	assert.Equal(NormalizeText(name), "BurgisweiherSprint125eae")
}
