// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha_memory_storage

import (
	"github.com/orivil/service"
	"github.com/orivil/services/captcha"
)

type Service struct {
	self service.Provider
}

func (s *Service) Get(ctn *service.Container) (captcha.Storage, error) {
	store, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return captcha.Storage(store.(*Storage)), nil
	}
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	return NewStorage(), nil
}

func NewService() *Service {
	s := &Service{}
	s.self = s
	return s
}
