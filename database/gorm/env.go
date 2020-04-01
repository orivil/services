// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package gorm

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"sync"
)

type DBDialect string

const (
	Mysql    DBDialect = "mysql"
	Postgres DBDialect = "postgres"
	SQLite3  DBDialect = "sqlite3"
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

type Env struct {
	Debug       bool   `toml:"debug"`
	Dialect     string `toml:"dialect"`
	TablePrefix string `toml:"table_prefix"`
}

var once = sync.Once{}

func (e *Env) Init(DB *sql.DB) (*gorm.DB, error) {
	db, err := gorm.Open(e.Dialect, DB)
	if err != nil {
		return nil, err
	}
	db.LogMode(e.Debug)
	if e.TablePrefix != "" {
		db.Set("table_prefix", e.TablePrefix) // 紧设置当前 db 的表前缀
		once.Do(func() {                      // 全局方法只需要设置一次
			gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
				if tf, ok := db.Get("table_prefix"); ok {
					return tf.(string) + defaultTableName
				}
				return defaultTableName
			}
		})
	}
	return db, nil
}
