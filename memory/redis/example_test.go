// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package redis_test

import (
	"fmt"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/cfg/storages/memory"
	"github.com/orivil/services/memory/redis"
	"time"
)

func ExampleNewService() {
	var config = `
# redis配置
[redis]
# 是否模拟客户端, 开启后将链接至虚拟客户端, 测试时使用
is_mock = true
# 地址
addr = ""
# 密码
password = ""
`
	cfgService := cfg.NewService(memory.NewService(config))
	redisService := redis.NewService(cfgService, 1)
	container := service.NewContainer()
	client, err := redisService.Get(container)
	if err != nil {
		panic(err)
	}
	client.Set("key", "value", time.Second*2)
	var v1 string
	v1, err = client.Get("key").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(v1) // value

	// 将模拟服务器时间向前推进
	redisService.GetMockServer().FastForward(2 * time.Second)
	var exist int64
	exist, err = client.Exists("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(exist == 1)
	// Output:
	// value
	// false
}
