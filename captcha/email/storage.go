// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email_captcha

import "io/ioutil"

type TemplateStorage interface {
	Read() ([]byte, error)
}

type TemplateFileStorage string

func (s TemplateFileStorage) Read() ([]byte, error) {
	return ioutil.ReadFile(string(s))
}

type TemplateMemoryStorage string

func (s TemplateMemoryStorage) Read() ([]byte, error) {
	return []byte(s), nil
}
