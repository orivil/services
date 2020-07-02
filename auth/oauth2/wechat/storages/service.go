// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package wechat_storages

import (
	redis2 "github.com/go-redis/redis/v8"
	"github.com/orivil/service"
	"github.com/orivil/services/auth/oauth2/wechat"
	"github.com/orivil/services/memory/redis"
)

type MemoryService struct {
	self service.Provider
}

func (m *MemoryService) New(ctn *service.Container) (value interface{}, err error) {
	return wechat.Storage(NewMemory()), nil
}

func (m *MemoryService) Get(ctn *service.Container) (wechat.Storage, error) {
	store, err := ctn.Get(&m.self)
	if err != nil {
		return nil, err
	} else {
		return store.(wechat.Storage), nil
	}
}

func NewMemoryService() *MemoryService {
	s := &MemoryService{}
	s.self = s
	return s
}

type RedisService struct {
	self  service.Provider
	user  *redis.Service
	token *redis.Service
}

func (r *RedisService) New(ctn *service.Container) (value interface{}, err error) {
	var (
		token *redis2.Client
		user  *redis2.Client
	)
	token, err = r.token.Get(ctn)
	if err != nil {
		return nil, err
	}
	user, err = r.user.Get(ctn)
	if err != nil {
		return nil, err
	}
	return wechat.Storage(NewRedis(token, user)), nil
}

func (r *RedisService) Get(ctn *service.Container) (wechat.Storage, error) {
	store, err := ctn.Get(&r.self)
	if err != nil {
		return nil, err
	} else {
		return store.(wechat.Storage), nil
	}
}

func NewRedisService(user, token *redis.Service) *RedisService {
	s := &RedisService{
		user:  user,
		token: token,
	}
	s.self = s
	return s
}
