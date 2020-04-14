// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email_captcha

import "github.com/orivil/services/captcha"

/**
# 邮箱验证码配置
[email-captcha]

# 邮件内容类型
content_type = "text/html; charset=UTF-8"

# 验证码过期时间(单位：秒)
expires = 300

# 字符集, 最好排除容易混淆的字符. 如果有小写字母，则将被自动转换为大写字母
chars = "2345789ABCDEFHJKLMNPRSTWXYZ"
*/

type Env struct {
	ContentType string `toml:"content_type"`
	captcha.Env
}
