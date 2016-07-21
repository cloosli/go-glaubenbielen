package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun(t *testing.T) {
	assert := assert.New(t)

	flagFilename = "../data/test"
	assert.Error(run())

	flagFilename = "../data/activity_1263337503.kml"
	assert.Error(run())

	flagFilename = "../data/activity_1263337503.gpx"
	assert.NoError(run())
}
