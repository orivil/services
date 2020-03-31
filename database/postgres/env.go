// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/orivil/services/database"
)

// postgres 配置参数
type Env struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
	SSLMode  string `toml:"ssl_mode"`
	database.Env
}

// DSN 数据库连接串
func (e Env) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		e.Host, e.Port, e.User, e.DBName, e.Password, e.SSLMode)
}

func (e Env) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", e.DSN())
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
