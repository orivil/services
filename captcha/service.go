// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

type Service struct {
	configService  service.Provider
	storageService service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	envs := ctn.MustGet(&s.configService).(xcfg.Env)
	env := &Env{}
	err = envs.UnmarshalSub("captcha", env)
	if err != nil {
		panic(err)
	}
	storage := ctn.MustGet(&s.storageService).(Storage)
	dispatcher := NewDispatcher(storage, env)
	return dispatcher, nil
}

func NewService(configService *cfg.Service, storageService service.Provider) *Service {
	return &Service{
		configService:  configService,
		storageService: storageService,
	}
}
