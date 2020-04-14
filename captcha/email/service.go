// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email_captcha

import (
	"github.com/orivil/service"
	"github.com/orivil/services/captcha"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/email"
	"github.com/orivil/xcfg"
	"text/template"
)

type Service struct {
	emailService    *email.Service
	configService   *cfg.Service
	configNamespace string
	storageService  captcha.StorageService
	ts              TemplateStorage
	self            service.Provider
}

func NewService(cfgNamespace string, cs *cfg.Service, es *email.Service, ss captcha.StorageService, ts TemplateStorage) *Service {
	s := &Service{
		emailService:    es,
		configService:   cs,
		configNamespace: cfgNamespace,
		storageService:  ss,
		ts:              ts,
		self:            nil,
	}
	s.self = s
	return s
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var (
		envs        xcfg.Env
		storage     captcha.Storage
		emailSender *email.Sender
		data        []byte
		tpl         *template.Template
	)
	envs, err = s.configService.Get(ctn)
	if err != nil {
		return nil, err
	}
	env := &Env{}
	err = envs.UnmarshalSub(s.configNamespace, env)
	if err != nil {
		return nil, err
	}
	storage, err = s.storageService.Get(ctn)
	if err != nil {
		return nil, err
	}
	emailSender, err = s.emailService.Get(ctn)
	if err != nil {
		return nil, err
	}
	data, err = s.ts.Read()
	if err != nil {
		return nil, err
	}
	tpl, err = template.New("email").Parse(string(data))
	if err != nil {
		return nil, err
	}
	return NewDispatcher(env, emailSender, storage, tpl), nil
}

func (s *Service) Get(ctn *service.Container) (*Dispatcher, error) {
	d, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	}
	return d.(*Dispatcher), nil
}
