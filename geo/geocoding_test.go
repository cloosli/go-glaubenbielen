package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLabel(t *testing.T) {
	// Skip tests that actually hit the API in short mode. We run with the
	// -test.short flag in Travis from external pull requests because the
	// credentials are not available.
	if testing.Short() {
		return
	}
	assert := assert.New(t)

	//home := os.Getenv("HOME")
	//filepath.Join(home, "data/campaigns")

	var filename string
	filename = "../data/test"
	assert.Error(run(filename))

	filename = "../data/activity_1263337503.kml"
	assert.Error(run(filename))

	filename = "../data/activity_1263337503.gpx"
	assert.NoError(run(filename))
}
