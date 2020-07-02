// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha_redis_storage

import (
	"github.com/go-redis/redis/v8"
	"github.com/orivil/service"
	"github.com/orivil/services/captcha"
	redis2 "github.com/orivil/services/memory/redis"
)

type Service struct {
	redisService *redis2.Service
	self         service.Provider
}

func (s *Service) Get(ctn *service.Container) (captcha.Storage, error) {
	store, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		storeInt := store.(*Storage)
		return captcha.Storage(storeInt), nil
	}
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var client *redis.Client
	client, err = s.redisService.Get(ctn)
	if err != nil {
		return nil, err
	}
	return NewStorage(client), nil
}

func NewService(redisService *redis2.Service) *Service {
	s := &Service{redisService: redisService}
	s.self = s
	return s
}
