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

/**
# gorm 配置
[gorm]
# 是否开启调试模式
debug = true
# 数据库类型(目前支持的数据库类型：mysql/sqlite3/postgres)
dialect = "sqlite3"
# 数据库表名前缀
table_prefix = ""
*/
var Service service.Provider = func(c *service.Container) (value interface{}, err error) {
	env := c.MustGet(&cfg.Service).(xcfg.Env)
	gormEnv := &Env{}
	err = env.UnmarshalSub("gorm", gormEnv)
	if err != nil {
		return nil, err
	} else {
		var db *sql.DB
		switch DBDialect(gormEnv.Dialect) {
		case Mysql:
			db = c.MustGet(&mysql.Service).(*sql.DB)
		case Postgres:
			db = c.MustGet(&postgres.Service).(*sql.DB)
		case SQLite3:
			db = c.MustGet(&sqlite.Service).(*sql.DB)
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
