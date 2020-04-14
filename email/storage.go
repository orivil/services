// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email

import "github.com/orivil/service"

type Storage interface {
	GetEnv() (*Env, error)
	SetEnv(env *Env) error
}

type StorageService interface {
	service.Provider
	Get(container *service.Container) (Storage, error)
}
