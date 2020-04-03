// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
)

type Service struct {
	cfgService      *cfg.Service
	templateStorage TemplateStorage
	cfgNamespace    string
	self            service.Provider
}

func (s *Service) New(container *service.Container) (interface{}, error) {
	envs, err := s.cfgService.Get(container)
	if err != nil {
		return nil, err
	}
	env := &Env{}
	err = envs.UnmarshalSub(s.cfgNamespace, env)
	if err != nil {
		return nil, err
	}
	return NewSMTPSender(env, s.templateStorage)
}

func (s *Service) Get(container *service.Container) (*Sender, error) {
	sender, err := container.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return sender.(*Sender), nil
	}
}

func NewService(configNamespace string, cfgService *cfg.Service, templateStorage TemplateStorage) *Service {
	s := &Service{
		cfgService:      cfgService,
		templateStorage: templateStorage,
		cfgNamespace:    configNamespace,
		self:            nil,
	}
	s.self = s
	return s
}
