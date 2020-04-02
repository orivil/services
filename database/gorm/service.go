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
	"github.com/orivil/services/database/mysql"
	"github.com/orivil/services/database/postgres"
	"github.com/orivil/services/database/sqlite"
	"github.com/orivil/xcfg"
)

type Service struct {
	configService   *cfg.Service
	mysqlService    *mysql.Service
	postgresService *postgres.Service
	sqliteService   *sqlite.Service
	self            service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var envs xcfg.Env
	envs, err = s.configService.Get(ctn)
	if err != nil {
		return nil, err
	}
	env := &Env{}
	err = envs.UnmarshalSub("gorm", env)
	if err != nil {
		return nil, err
	} else {
		var db *sql.DB
		switch DBDialect(env.Dialect) {
		case Mysql:
			db, err = s.mysqlService.Get(ctn)
		case Postgres:
			db, err = s.postgresService.Get(ctn)
		case SQLite3:
			db, err = s.sqliteService.Get(ctn)
		default:
			return nil, fmt.Errorf("gorm: dialect [%s] is not supported", env.Dialect)
		}
		var gormDB, er = env.Init(db)
		if er != nil {
			return nil, er
		} else {
			return gormDB, nil
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

// 新建 gorm 服务，参数 sqlService 只能是 *mysql.Service, *postgres.Service 或者 *sqlite.Service
func NewService(configService *cfg.Service, sqlService ...interface{}) *Service {
	s := &Service{
		configService:   configService,
		mysqlService:    nil,
		postgresService: nil,
		sqliteService:   nil,
	}
	for _, se := range sqlService {
		s.addSqlService(se)
	}
	s.self = s
	return s
}

func (s *Service) addSqlService(service interface{}) {
	switch v := service.(type) {
	case *mysql.Service:
		s.mysqlService = v
	case *postgres.Service:
		s.postgresService = v
	case *sqlite.Service:
		s.sqliteService = v
	default:
		panic("sql service not allowed")
	}
}
