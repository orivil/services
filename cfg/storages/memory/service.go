// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package memory

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
)

type Service struct {
	store Storage
	self  service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	return cfg.Storage(s.store), nil
}

func (s *Service) Get(ctn *service.Container) (cfg.Storage, error) {
	store, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return store.(cfg.Storage), nil
	}
}

func NewService(configData string) *Service {
	s := &Service{
		store: Storage(configData),
		self:  nil,
	}
	s.self = s
	return s
}
