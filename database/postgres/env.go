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

/**
# postgres 数据库配置
[postgres]
# 连接地址
host = "127.0.0.1"
# 连接端口
port= 5432
# 用户名
user = "root"
# 密码
password = "123456"
# 数据库
db_name = "ginadmin"
# SSL模式
ssl_mode = "disable"
# 设置连接可以重用的最长时间(单位：秒)
max_lifetime = 7200
# 设置数据库的最大打开连接数
max_open_conns = 150
# 设置空闲连接池中的最大连接数
max_idle_conns = 50
*/
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
