// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha

import "crypto/rand"

// ID 字符集
var chars = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 从 src 中随机创建 bit 位的随机字符数组
func RandomBytes(src []byte, bit int) []byte {
	bts := make([]byte, bit)
	r := make([]byte, bit)
	ln := len(src)
	rand.Read(bts)
	for idx, b := range bts {
		r[idx] = src[int(b)%ln]
	}
	return r
}

func NewID(bit int) []byte {
	return RandomBytes(chars, bit)
}
