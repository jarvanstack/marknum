package main

import (
	"fmt"
	"testing"
)

func Test_mdPaths(t *testing.T) {
	paths := mdPaths("./testdata")
	fmt.Printf("paths=%#v\n", paths)
}
