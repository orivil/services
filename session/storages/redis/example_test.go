// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session_redis_storage

import (
	"fmt"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/memory/redis"
	"time"
)

// redis 存储器测试, redis 服务器为模拟服务器, 与实际 redis 服务器可能有所区别
// see: github.com/alicebob/miniredis/v2
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
# 数据库
db = 1
`
	cfgService := cfg.NewService(cfg.NewMemoryStorageService(config))
	redisService := redis.NewService("redis", cfgService)
	storeService := NewService(redisService)
	container := service.NewContainer()
	store, _ := storeService.Get(container)
	session := "tokenStr"

	err := store.SaveSession(session, 2*time.Second)
	panicErr(err)
	var ok bool
	ok, err = store.IsOK(session)
	panicErr(err)
	fmt.Println(ok == true)

	// 将 redis 模拟服务器的时间向前推进, 以删除过期数据
	redisService.GetMockServer().FastForward(4 * time.Second)

	// 已过期
	ok, err = store.IsOK(session)
	panicErr(err)
	fmt.Println(ok == false)

	err = store.SaveSession(session, 5*time.Minute)
	panicErr(err)
	ok, err = store.IsOK(session)
	panicErr(err)
	fmt.Println(ok == true)
	err = store.DelSession(session)
	panicErr(err)
	ok, err = store.IsOK(session)
	panicErr(err)
	fmt.Println(ok == false)

	// Output:
	// true
	// true
	// true
	// true
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
