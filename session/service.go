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
	configService  service.Provider
	storageService service.Provider
}

func (J *JWTAuthService) New(ctn *service.Container) (value interface{}, err error) {
	envs := ctn.MustGet(&J.configService).(xcfg.Env)
	env := JWTEnv{}
	err = envs.UnmarshalSub("jwt_auth", &env)
	if err != nil {
		panic(err)
	}
	storage := ctn.MustGet(&J.storageService).(Storage)
	return NewJWTAuth(storage, env), nil
}

// 参数 storageService 参考 redis_storage 目录下的 redis_storage.Service
func NewJWTAuthService(configService *cfg.Service, storageService service.Provider) *JWTAuthService {
	return &JWTAuthService{
		configService:  configService,
		storageService: storageService,
	}
}
