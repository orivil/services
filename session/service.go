// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

type Service struct {
	configService   *cfg.Service
	storageService  StorageService
	configNamespace string
	self            service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var envs xcfg.Env
	envs, err = s.configService.Get(ctn)
	if err != nil {
		return nil, err
	}
	env := Env{}
	err = envs.UnmarshalSub(s.configNamespace, &env)
	if err != nil {
		return nil, err
	}
	var store Storage
	store, err = s.storageService.Get(ctn)
	if err != nil {
		return nil, err
	}
	return NewDispatcher(store, env), nil
}

func (s *Service) Get(ctn *service.Container) (*Dispatcher, error) {
	auth, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return auth.(*Dispatcher), nil
	}
}

func NewService(configNamespace string, configService *cfg.Service, storageService StorageService) *Service {
	s := &Service{
		configService:   configService,
		configNamespace: configNamespace,
		storageService:  storageService,
	}
	s.self = s
	return s
}
