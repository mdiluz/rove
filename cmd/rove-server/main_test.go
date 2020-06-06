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

func Test_InnerMain_Quit(t *testing.T) {
	flag.Set("quit", "1")
	InnerMain()
	flag.Set("quit", "0")
}
