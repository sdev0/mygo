//go:build ignore
// +build ignore

package main

import (
	_ "github.com/sdev0/mygo/apis"
	"github.com/sdev0/mygo/test"
)

func main() {
	// apis.InitAll()
	// apis.Spider92qb()
	test.DoTest()
}
