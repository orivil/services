// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package limiter

import (
	"github.com/orivil/limiter"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

type Service struct {
	cfgService *cfg.Service
	storage    StorageService
	namespace  string
	self       service.Provider
}

func NewService(namespace string, cfgService *cfg.Service, storageService StorageService) *Service {
	s := &Service{cfgService: cfgService, namespace: namespace, storage: storageService}
	s.self = s
	return s
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var envs xcfg.Env
	envs, err = s.cfgService.Get(ctn)
	if err != nil {
		return nil, err
	}
	env := &Env{}
	err = envs.UnmarshalSub(s.namespace, env)
	if err != nil {
		return nil, err
	}
	opts := toLimiterOptions(env)
	var storage Storage
	storage, err = s.storage.Get(ctn)
	return limiter.NewTimesLimiter(opts, storage), nil
}

func (s *Service) Get(ctn *service.Container) (*limiter.TimesLimiter, error) {
	l, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	}
	return l.(*limiter.TimesLimiter), nil
}
