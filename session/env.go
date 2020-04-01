// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

/**
# 用户认证(jwt)
[jwt_auth]
# 签名方式(支持：HS512/HS384/HS256)
signing_method = "HS512"
# 签名key
signing_key = "GINADMIN"
# 过期时间（单位秒）
expired = 7200
# 过期前刷新时间（单位秒），如值为 1800 则在过期前 30 分钟刷新
refresh = 1800
# 存储(支持：file/redis)
store = "file"
# 文件路径
file_path = "data/jwt_auth.db"
# redis数据库(如果存储方式是redis，则指定存储的数据库)
redis_db = 10
# 存储到redis数据库中的键名前缀
redis_prefix = "auth_"
*/
type JWTEnv struct {
	SigningMethod string `toml:"signing_method"`
	SigningKey    string `toml:"signing_key"`
	Expired       int64  `toml:"expired"`
	Refresh       int64  `toml:"refresh"`
	Store         string `toml:"store"`
	FilePath      string `toml:"file_path"`
	RedisDB       string `toml:"redis_db"`
	RedisPrefix   string `toml:"redis_prefix"`
}
