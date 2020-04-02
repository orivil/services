// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package local

import "io/ioutil"

type FileStorage string

func (f FileStorage) GetTomlData() ([]byte, error) {
	return ioutil.ReadFile(string(f))
}
