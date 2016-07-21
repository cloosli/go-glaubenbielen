package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFilename(t *testing.T) {
	assert := assert.New(t)
	name := "Bürgisweiher Sprint 125 éàè!!"
	fmt.Println(strings.Title(name))
	fmt.Println(strings.ToTitle(name))
	fmt.Println(NormalizeText(name))
	assert.Equal(NormalizeText(name), "BurgisweiherSprint125eae")
}

func TestRun(t *testing.T) {
	assert := assert.New(t)

	var filename string
	filename = "../data/test"
	assert.Error(run(filename))

	filename = "../data/activity_1263337503.kml"
	assert.Error(run(filename))

	filename = "../data/activity_1263337503.gpx"
	assert.NoError(run(filename))
}
