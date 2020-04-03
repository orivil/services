// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

///+build ignore

package main

import (
	"github.com/orivil/service"
	"github.com/orivil/services/auth"
	"github.com/orivil/services/auth/storages/gorm"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/cfg/storages/config_memory_storage"
	"github.com/orivil/services/database/gorm"
	"github.com/orivil/services/database/sqlite"
	"github.com/orivil/services/memory/redis"
	"github.com/orivil/services/session"
	"github.com/orivil/services/session/storages/redis"
)


var emailTemplate = `<!doctype html>
<html lang="en">
<head>
   <meta charset="UTF-8">
   <title>{{.title}}</title>
</head>
<body>
<h1 style="text-align: center">{{.title}}</h1>
<p><span style="font-weight: bold">{{.code}}</span> {{.message}}</p>
</body>
</html>`

var config = `
# sqlite数据库配置
[sqlite3]
# 数据库路径
path = "data/SQLite.db"

# gorm 配置
[gorm]
# 是否开启调试模式
debug = true
# 数据库类型(目前支持的数据库类型：mysql/sqlite3/postgres)
dialect = "sqlite3"
# 数据库表名前缀
table_prefix = ""

# redis配置
[redis]
# 是否模拟客户端, 开启后将链接至虚拟客户端, 测试时使用
# 使用方式参考: https://github.com/alicebob/miniredis
is_mock = true
# 地址
addr = ""
# 密码
password = ""
# 数据库
db = 1

# 用户会话(jwt认证)
[sessions]
# jwt 签名方式(支持：HS512/HS384/HS256)
signing_method = "HS512"
# jwt 签名key
signing_key = "GINADMIN"
# session 过期时间（单位秒）
expired = 7200
# session 过期之前的刷新时间(单位秒), 如果该不大于0, 则不刷新
refresh = 1800
`

// 测试该项目需要安装 SQLite3, 或则更改为 Postgres/Mysql
func main() {
	cfgService := cfg.NewService(cfg.NewMemoryStorageService(config))

	// gorm
	SQLiteService := sqlite.NewService("sqlite3", cfgService)
	gormService := gorm.NewService("gorm", cfgService, SQLiteService)

	redisService := redis.NewService("redis", cfgService)

	authGormStorageService := auth_gorm_storage.NewService(gormService)

	sessionRedisStorageService := session_redis_storage.NewService(redisService)

	sessionService := session.NewService("sessions", cfgService, sessionRedisStorageService)

	authDispatcherService := auth.NewService(sessionService, authGormStorageService)

	container := service.NewContainer()

	authDispatcher, err := authDispatcherService.Get(container)
	if err != nil {
		panic(err)
	}

	authDispatcher.
}
