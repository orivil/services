// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package wechat

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

type Service struct {
	cfg *cfg.Service
	cfgNamespace string
	store  StorageService
	self service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var envs xcfg.Env
	envs, err = s.cfg.Get(ctn)
	if err != nil {
		return nil, err
	}
	env := &Config{}
	err = envs.UnmarshalSub(s.cfgNamespace, env)
	if err != nil {
		return nil, err
	}
	var store Storage
	store, err = s.store.Get(ctn)
	if err != nil {
		return nil, err
	}
	return NewDispatcher(env, store), nil
}

func (s *Service) Get(ctn *service.Container) (*Dispatcher, error) {
	d, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return d.(*Dispatcher), nil
	}
}

type Options struct {
	ConfigNamespace string
	ConfigService *cfg.Service
	Storage StorageService
}

func NewService(opts Options) *Service {
	s := &Service {
		cfg: opts.ConfigService,
		cfgNamespace: opts.ConfigNamespace,
		store: opts.Storage,
	}
	s.self = s
	return s
}