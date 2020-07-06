package main

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InnerMain_Version(t *testing.T) {
	assert.NoError(t, flag.Set("version", "1"))
	InnerMain()
	assert.NoError(t, flag.Set("version", "0"))
}
