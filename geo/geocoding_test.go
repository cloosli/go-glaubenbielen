package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"

	"github.com/cloosli/go-glaubenbielen/util"
)

func TestFilename(t *testing.T) {
	assert := assert.New(t)
	name := "Bürgisweiher Sprint 125 éàè!!"
	fmt.Println(strings.Title(name))
	fmt.Println(strings.ToTitle(name))
	fmt.Println(util.NormalizeText(name))
	assert.Equal(util.NormalizeText(name), "BurgisweiherSprint125eae")
}

func TestRun(t *testing.T) {
	assert := assert.New(t)

	flagFilename = "../data/test"
	assert.Error(run())

	flagFilename = "../data/activity_1263337503.kml"
	assert.Error(run())

	flagFilename = "../data/activity_1263337503.gpx"
	assert.NoError(run())
}
