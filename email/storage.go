// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email

import "io/ioutil"

type TemplateStorage interface {
	Read() ([]byte, error)
}

type FileStorage string

func (s FileStorage) Read() ([]byte, error) {
	return ioutil.ReadFile(string(s))
}

type MemoryStorage string

func (s MemoryStorage) Read() ([]byte, error) {
	return []byte(s), nil
}
