// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha

import "strings"

/**
# 验证码长度
captcha_length = 4

# 验证码过期时间(单位：秒)
expires = 300

# 字符集, 最好排除容易混淆的字符. 如果有小写字母，则将被自动转换为大写字母
chars = "2345789ABCDEFHJKLMNPRSTWXYZ"
*/
type Env struct {
	CaptchaLength int    `toml:"captcha_length"`
	Expires       int64  `toml:"expires"`
	Chars         string `toml:"chars"`
}

func (e *Env) GenCaptcha() string {
	return strings.ToUpper(string(RandomBytes([]byte(e.Chars), e.CaptchaLength)))
}
