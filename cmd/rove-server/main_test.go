package main

import (
	"flag"
	"testing"
)

func Test_InnerMain_Version(t *testing.T) {
	flag.Set("version", "1")
	InnerMain()
	flag.Set("version", "0")
}
