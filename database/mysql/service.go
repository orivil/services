// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package mysql

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

/**
# mysql数据库配置
[mysql]
# 连接地址
host = "127.0.0.1"
# 连接端口
port= "3306"
# 用户名
user = "root"
# 密码
password = "123456"
# 数据库
db_name = ""
# 连接参数
parameters = "charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true"
# 设置连接可以重用的最长时间(单位：秒)
max_lifetime = 7200
# 设置数据库的最大打开连接数
max_open_conns = 150
# 设置空闲连接池中的最大连接数
max_idle_conns = 50
*/

var Service service.Provider = func(c *service.Container) (value interface{}, err error) {
	env := c.MustGet(&cfg.Service).(xcfg.Env)
	mysql := &Env{}
	err = env.UnmarshalSub("mysql", mysql)
	if err != nil {
		return nil, err
	}
	return mysql.Connect()
}
