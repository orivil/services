// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

/**
# 用户会话(jwt认证)
[sessions]
# jwt 签名方式(支持：HS512/HS384/HS256)
signing_method = "HS512"
# jwt 签名key
signing_key = "secret key"
# session 过期时间（单位秒）
expires = 7200
*/
type Env struct {
	SigningMethod string `toml:"signing_method"`
	SigningKey    string `toml:"signing_key"`
	Expires       int64  `toml:"expires"`
}
