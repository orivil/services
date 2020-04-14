// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email

/**
# SMTP 邮箱配置
[smtp-email]

# 域名
host = "smtp.gmail.com"

# 端口
port = ":587"

# 用户名
username = "xxx"

# 密码
password = "xxx"

# 来源
from = "orivil.com"
*/
type Env struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	From     string `toml:"from"`
}
