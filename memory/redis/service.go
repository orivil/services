// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package redis

import (
	"github.com/go-redis/redis"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

/**
# redis配置
[redis]
# 地址
addr = "127.0.0.1:6379"
# 密码
password = ""
*/

// 新建 service 需要在全局保存起来
func NewService(db int) service.Provider {
	return func(c *service.Container) (value interface{}, err error) {
		envs := c.MustGet(&cfg.Service).(xcfg.Env)
		redisEnv := &Env{}
		err = envs.UnmarshalSub("redis", redisEnv)
		if err != nil {
			return nil, err
		}
		var client *redis.Client
		client, err = redisEnv.Init(db)
		return client, err
	}
}
