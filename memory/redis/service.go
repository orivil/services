// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package redis

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

type Service struct {
	configService *cfg.Service
	db            int
	self          service.Provider
	mockServer    *miniredis.Miniredis
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var envs xcfg.Env
	envs, err = s.configService.Get(ctn)
	if err != nil {
		return nil, err
	}
	env := &Env{}
	err = envs.UnmarshalSub("redis", env)
	if err != nil {
		panic(err)
	}
	client, mockServer, err1 := env.Init(s.db)
	if err1 != nil {
		return nil, err1
	}
	s.mockServer = mockServer
	ctn.OnClose(func() error {
		client.Close()
		if mockServer != nil {
			mockServer.Close()
		}
		return nil
	})
	return client, nil
}

func (s *Service) Get(ctn *service.Container) (*redis.Client, error) {
	c, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return c.(*redis.Client), nil
	}
}

// 仅当配置参数 is_mock 为 true 时才能获得该值, 否则为 nil
func (s *Service) GetMockServer() *miniredis.Miniredis {
	return s.mockServer
}

func NewService(configService *cfg.Service, DB int) *Service {
	s := &Service{
		configService: configService,
		db:            DB,
	}
	s.self = s
	return s
}
