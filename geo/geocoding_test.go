package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLabel(t *testing.T) {
	assert := assert.New(t)

	var filename string
	filename = "../data/test"
	assert.Error(run(filename))

	filename = "../data/activity_1263337503.kml"
	assert.Error(run(filename))

	filename = "../data/activity_1263337503.gpx"
	assert.NoError(run(filename))
}
