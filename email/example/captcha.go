// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/email"
)

var config = `
# SMTP 邮箱配置
[smtp-email]

# 域名
host = "smtp.qq.com"

# 端口
port = ":25"

# 用户名
username = "xxx@xx.com"

# 密码
password = "xxx"

# 来源
from = "xxx@xx.com"`

func main() {
	cfgService := cfg.NewService(cfg.NewMemoryStorageService(config))
	emailService := email.NewService("smtp-email", cfgService)
	container := service.NewContainer()
	sender, err := emailService.Get(container)
	if err != nil {
		panic(err)
	}
	err = sender.Send([]string{"ksnmmmm@qq.com"}, "verify email", "text/html; charset=UTF-8", []byte(template))
	if err != nil {
		panic(err)
	}
}

var template = `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.title}}</title>
</head>
<body>
<h1 style="text-align: center">Verify your Email</h1>
<p><span style="font-weight: bold">XFK9</span> is your captcha</p>
</body>
</html>`
