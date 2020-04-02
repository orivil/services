// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package sqlite

import (
	"database/sql"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

type Service struct {
	configService *cfg.Service
	self          service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var envs xcfg.Env
	envs, err = s.configService.Get(ctn)
	if err != nil {
		return nil, err
	}
	env := &Env{}
	err = envs.UnmarshalSub("sqlite3", env)
	if err != nil {
		panic(err)
	}
	db, er := env.Connect()
	if er != nil {
		panic(er)
	}
	ctn.OnClose(func() error {
		return db.Close()
	})
	return db, nil
}

func (s *Service) Get(ctn *service.Container) (db *sql.DB, err error) {
	d, er := ctn.Get(&s.self)
	if er != nil {
		return nil, er
	} else {
		return d.(*sql.DB), nil
	}
}

func NewService(configService *cfg.Service) *Service {
	s := &Service{configService: configService}
	s.self = s
	return s
}
