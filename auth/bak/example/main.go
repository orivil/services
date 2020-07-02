// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

///+build ignore

package main

import (
	"github.com/orivil/service"
	"github.com/orivil/services/auth/bak"
	"github.com/orivil/services/auth/bak/storages/gorm"
	email2 "github.com/orivil/services/captcha/email"
	"github.com/orivil/services/captcha/image"
	"github.com/orivil/services/captcha/storages/memory"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/database/gorm"
	"github.com/orivil/services/database/sqlite"
	"github.com/orivil/services/email"
	"github.com/orivil/services/limiter"
	rds "github.com/orivil/services/memory/redis"
	"github.com/orivil/services/session"
	"github.com/orivil/services/session/storages/redis"
	"net/http"
	"path"
	"path/filepath"
)

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
	var (
		cfgService            *cfg.Service
		SQLiteService         *sqlite.Service
		gormService           *gorm.Service
		sessionStorageService session.StorageService
		sessionService        *session.Service
		imageCaptchaService   *image_captcha.Service
		emailService          *email.Service
		emailCaptchaService   *email2.Service
		limiterService        *limiter.Service
		authStorageService    bak.StorageService
		authDispatcherService *bak.Service
	)
	cfgService = cfg.NewService(cfg.NewMemoryStorageService(config))

	SQLiteService = sqlite.NewService("sqlite3", cfgService)

	gormService = gorm.NewService("gorm", cfgService, SQLiteService)

	sessionStorageService = session_redis_storage.NewService(rds.NewService("session-redis", cfgService))

	sessionService = session.NewService("sessions", cfgService, sessionStorageService)

	imageCaptchaService = image_captcha.NewService("image-captcha", cfgService, captcha_memory_storage.NewService())

	emailService = email.NewService("smtp-email", cfgService)

	emailCaptchaService = email2.NewService("email-captcha", cfgService, emailService, captcha_memory_storage.NewService(), email2.TemplateMemoryStorage(emailTemplate))

	limiterService = limiter.NewService("email-limiter", cfgService, limiter.NewMemoryStorageService())

	authStorageService = auth_gorm_storage.NewService(gormService)

	authDispatcherService = bak.NewService(authStorageService, sessionService, imageCaptchaService, emailCaptchaService, limiterService)

	container := service.NewContainer()

	dispatcher, err := authDispatcherService.Get(container)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, file := path.Split(request.URL.Path)
		if file == "" {
			file = "index"
		}
		http.ServeFile(writer, request, filepath.Join("pages", file+".gohtml"))
	})
}

var emailTemplate = `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Verify Email</title>
</head>
<body>
<h1 style="text-align: center">Verify your Email</h1>
<p><span style="font-weight: bold">{{.}}</span> is your captcha</p>
</body>
</html>`

var indexPage = []byte(``)
