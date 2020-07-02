// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package gorm

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/database"
	"github.com/orivil/xcfg"
)

type DatabaseService interface {
	Get(ctn *service.Container) (db *sql.DB, err error)
	Dialect() database.DBDialect
}

type Service struct {
	configService   *cfg.Service
	dbServices      map[string]DatabaseService
	configNamespace string
	self            service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var envs xcfg.Env
	envs, err = s.configService.Get(ctn)
	if err != nil {
		return nil, err
	}
	env := &Env{}
	err = envs.UnmarshalSub(s.configNamespace, env)
	if err != nil {
		return nil, err
	} else {
		var db *sql.DB
		if dbService, ok := s.dbServices[env.Dialect]; ok {
			db, err = dbService.Get(ctn)
			if err != nil {
				return nil, err
			}
		} else {
			panic(fmt.Sprintf("database service [%s] not provided", env.Dialect))
		}
		var gDB *gorm.DB
		gDB, err = env.Init(db)
		if err != nil {
			return nil, err
		} else {
			return gDB, nil
		}
	}
}

func (s *Service) Get(ctn *service.Container) (*gorm.DB, error) {
	db, er := ctn.Get(&s.self)
	if er != nil {
		return nil, er
	} else {
		return db.(*gorm.DB), nil
	}
}

// 新建 gorm 服务
func NewService(configNamespace string, configService *cfg.Service, dbService ...DatabaseService) *Service {
	s := &Service{
		configService:   configService,
		configNamespace: configNamespace,
		dbServices: func() map[string]DatabaseService {
			mp := make(map[string]DatabaseService, len(dbService))
			for _, databaseService := range dbService {
				mp[string(databaseService.Dialect())] = databaseService
			}
			return mp
		}(),
	}
	s.self = s
	return s
}
