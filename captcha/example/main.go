// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

// +build ignore

package main

import (
	"github.com/orivil/services/captcha"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		code := request.URL.Query().Get("code")
		img := captcha.NewImage("1123", []byte(code), 280, 80)
		_, err := img.WriteTo(writer)
		if err != nil {
			panic(err)
		}
	})
	err := http.ListenAndServe(":8882", nil)
	panic(err)
}
