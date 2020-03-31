// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package postgres

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
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

var Service service.Provider = func(c *service.Container) (value interface{}, err error) {
	env := c.MustGet(&cfg.Service).(xcfg.Env)
	postgres := &Env{}
	err = env.UnmarshalSub("postgres", postgres)
	if err != nil {
		return nil, err
	}
	var db, er = postgres.Connect()
	if er != nil {
		return nil, er
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
