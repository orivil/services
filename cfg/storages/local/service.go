// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package local

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
)

const DefaultConfigFile = "configs/config.toml"

type Service struct {
	self     service.Provider
	filename string
}

func (s *Service) Get(ctn *service.Container) (cfg.Storage, error) {
	store, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return store.(cfg.Storage), nil
	}
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	return cfg.Storage(FileStorage(s.filename)), nil
}

func NewService(filename string) *Service {
	if filename == "" {
		filename = DefaultConfigFile
	}
	s := &Service{
		self:     nil,
		filename: filename,
	}
	s.self = s
	return s
}
