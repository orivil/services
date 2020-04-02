// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import "github.com/orivil/service"

type Storage interface {
	GetTomlData() ([]byte, error)
}

type StorageService interface {
	service.Provider
	Get(ctn *service.Container) (Storage, error)
}
