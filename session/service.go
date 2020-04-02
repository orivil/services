// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

type JWTAuthService struct {
	configService  *cfg.Service
	storageService StorageService
	self           service.Provider
}

func (j *JWTAuthService) New(ctn *service.Container) (value interface{}, err error) {
	var envs xcfg.Env
	envs, err = j.configService.Get(ctn)
	if err != nil {
		return nil, err
	}
	env := JWTEnv{}
	err = envs.UnmarshalSub("jwt_auth", &env)
	if err != nil {
		return nil, err
	}
	var store Storage
	store, err = j.storageService.Get(ctn)
	if err != nil {
		return nil, err
	}
	return NewJWTAuth(store, env), nil
}

func (j *JWTAuthService) Get(ctn *service.Container) (*JWTAuth, error) {
	auth, err := ctn.Get(&j.self)
	if err != nil {
		return nil, err
	} else {
		return auth.(*JWTAuth), nil
	}
}

func NewJWTAuthService(configService *cfg.Service, storageService StorageService) *JWTAuthService {
	s := &JWTAuthService{
		configService:  configService,
		storageService: storageService,
	}
	s.self = s
	return s
}
