// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package auth

import "github.com/orivil/service"

type StorageService func() (Storage, error)

func (sp StorageService) Init() service.Provider {
	return func(c *service.Container) (value interface{}, err error) {
		return sp()
	}
}

func NewService(storageService StorageService) ser
