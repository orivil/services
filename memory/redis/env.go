// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package redis

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
)

/**
# redis配置
[redis]
# 是否模拟客户端, 开启后将链接至虚拟客户端, 测试时使用
# 使用方式参考: https://github.com/alicebob/miniredis
is_mock = false
# 地址
addr = "127.0.0.1:6379"
# 密码
password = ""
*/

type Env struct {
	IsMock   bool   `toml:"is_mock"`
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

func (e Env) Init(db int) (client *redis.Client, mockServer *miniredis.Miniredis, err error) {
	var addr string
	if e.IsMock {
		mockServer, err = miniredis.Run()
		if err != nil {
			return nil, nil, err
		}

		addr = mockServer.Addr()
	} else {
		addr = e.Addr
	}
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: e.Password,
		DB:       db,
	})
	err = client.Ping().Err()
	if err != nil {
		return nil, nil, fmt.Errorf("ping redis: %s", err)
	} else {
		return client, mockServer, nil
	}
}
