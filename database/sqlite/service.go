// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package sqlite

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

/**
# sqlite数据库配置
[sqlite3]
# 数据库路径
path = "data/sqlite.db"
# 设置连接可以重用的最长时间(单位：秒)
max_lifetime = 7200
# 设置数据库的最大打开连接数
max_open_conns = 150
# 设置空闲连接池中的最大连接数
max_idle_conns = 50
*/

var Service service.Provider = func(c *service.Container) (value interface{}, err error) {
	env := c.MustGet(&cfg.Service).(xcfg.Env)
	sqlite := &Env{}
	err = env.UnmarshalSub("sqlite3", sqlite)
	if err != nil {
		return nil, err
	}
	var db, er = sqlite.Connect()
	if er != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
