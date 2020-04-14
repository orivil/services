// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"fmt"
	"path"
)

const (
	number = 1 << iota
	letter
)

func main() {
	dir, file := path.Split("/")
	fmt.Println("dir", dir, "file", file)
}
