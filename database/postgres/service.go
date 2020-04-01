// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package postgres

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

type Service struct {
	conf service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	envs := ctn.MustGet(&s.conf).(xcfg.Env)
	env := &Env{}
	err = envs.UnmarshalSub("postgres", env)
	if err != nil {
		panic(err)
	}
	db, er := env.Connect()
	if er != nil {
		panic(er)
	}
	return db, nil
}

func NewService(configService *cfg.Service) *Service {
	return &Service{conf: service.Provider(configService)}
}
