// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/orivil/services/database"
)

// sqlite3 配置参数
type Env struct {
	Path string `toml:"path"`
	database.Env
}

// DSN 数据库连接串
func (e Env) DSN() string {
	return e.Path
}

func (e Env) Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", e.DSN())
	if err != nil {
		return nil, err
	}
	err = e.Env.Init(db)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
