// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package image_captcha

import "github.com/orivil/services/captcha"

/**
# 验证码配置
[image-captcha]

# 验证码图片宽度(单位: px)
img_width = 240

# 验证码图片高度(单位：px)
img_height = 80

# 验证码长度
captcha_length = 4

# 验证码过期时间(单位：秒)
expires = 300

# 字符集, 最好排除容易混淆的字符. 如果有小写字母，则将被自动转换为大写字母
chars = "2345789ABCDEFHJKLMNPRSTWXYZ"
*/

type Env struct {
	ImgWidth  int `toml:"img_width"`
	ImgHeight int `toml:"img_height"`
	captcha.Env
}
