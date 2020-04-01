// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha

const (
	Number captchaType = 1 << iota
	Letter
	LetterAndNumber
)

type captchaType int

/**
# 验证码配置
[captcha]
# 验证码图片宽度(单位: px)
img_width = 120
# 验证码图片高度(单位：px)
img_height = 40
# 验证码长度
code_length = 6
# 验证码过期时间(单位：秒)
expires = 300
# 验证码类型：1-数字，2-字母，3-数字和字母
type = 3
*/

type Env struct {
	ImgWidth   int         `toml:"img_width"`
	ImgHeight  int         `toml:"img_height"`
	CodeLength int         `toml:"code_length"`
	Expires    int64       `toml:"expires"`
	Type       captchaType `toml:"type"`
}
