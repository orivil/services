// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package redis

import (
	"github.com/go-redis/redis"
	"github.com/orivil/service"
	redis2 "github.com/orivil/services/memory/redis"
	"github.com/orivil/services/session"
)

type Service struct {
	redisService *redis2.Service
	self         service.Provider
}

func (s *Service) Get(ctn *service.Container) (session.Storage, error) {
	store, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return store.(session.Storage), nil
	}
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var client *redis.Client
	client, err = s.redisService.Get(ctn)
	if err != nil {
		return nil, err
	}
	return session.Storage(NewStorage(client)), nil
}

func NewService(redisService *redis2.Service) *Service {
	s := &Service{
		redisService: redisService,
	}
	s.self = s
	return s
}
