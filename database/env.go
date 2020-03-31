// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package database

import (
	"database/sql"
	"time"
)

type Env struct {
	MaxLifetime  int `toml:"max_lifetime"`
	MaxOpenConns int `toml:"max_open_conns"`
	MaxIdleConns int `toml:"max_idle_conns"`
}

func (e Env) Init(db *sql.DB) error {
	if e.MaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(e.MaxLifetime) * time.Second)
	}
	if e.MaxIdleConns > 0 {
		db.SetMaxIdleConns(e.MaxIdleConns)
	}
	if e.MaxOpenConns > 0 {
		db.SetMaxOpenConns(e.MaxOpenConns)
	}
	return db.Ping()
}
