// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package gorm

import (
	"database/sql"
	"fmt"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/database/mysql"
	"github.com/orivil/services/database/postgres"
	"github.com/orivil/services/database/sqlite"
	"github.com/orivil/xcfg"
)

type Service struct {
	confService     service.Provider
	mysqlService    service.Provider
	postgresService service.Provider
	sqliteService   service.Provider
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

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	envs := ctn.MustGet(&s.confService).(xcfg.Env)
	gormEnv := &Env{}
	err = envs.UnmarshalSub("gorm", gormEnv)
	if err != nil {
		return nil, err
	} else {
		var db *sql.DB
		switch DBDialect(gormEnv.Dialect) {
		case Mysql:
			db = ctn.MustGet(&s.mysqlService).(*sql.DB)
		case Postgres:
			db = ctn.MustGet(&s.postgresService).(*sql.DB)
		case SQLite3:
			db = ctn.MustGet(&s.sqliteService).(*sql.DB)
		default:
			return nil, fmt.Errorf("gorm: dialect [%s] is not supported", gormEnv.Dialect)
		}
		var gormDB, er = gormEnv.Init(db)
		if er != nil {
			return nil, er
		} else {
			return gormDB, nil
		}
	}
}

// 新建 gorm 服务，参数 sqlService 只能是 *mysql.Service, *postgres.Service 或者 *sqlite.Service
func NewService(configService *cfg.Service, sqlService ...interface{}) *Service {
	s := &Service{
		confService:     configService,
		mysqlService:    nil,
		postgresService: nil,
		sqliteService:   nil,
	}
	for _, se := range sqlService {
		s.addSqlService(se)
	}
	return s
}
