// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

/// +build ignore

package main

import (
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	rds "github.com/orivil/services/memory/redis"
	"github.com/orivil/services/session"
	session_redis_storage "github.com/orivil/services/session/storages/redis"
	"strconv"
	"strings"
	"time"
)

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

# 用户会话(jwt认证)
[sessions]
# jwt 签名方式(支持：HS512/HS384/HS256)
signing_method = "HS256"
# jwt 签名key
signing_key = "GINADMIN"
# session 过期时间（单位秒）
expired = 7200
`

var redisMockServer *miniredis.Miniredis

func main() {
	cfgService := cfg.NewService(cfg.NewMemoryStorageService(config))
	redisService := rds.NewService("redis", cfgService)
	storeService := session_redis_storage.NewService(redisService)
	sessionService := session.NewService("sessions", cfgService, storeService)
	container := service.NewContainer()
	sessionDispatcher, err := sessionService.Get(container)
	if err != nil {
		panic(err)
	}
	// 只有调用 Get 方法之后才会注入依赖, 如果放在前面则获取不到 mock server
	redisMockServer = redisService.GetMockServer()

	testUnmarshal(sessionDispatcher)
	testExpired(sessionDispatcher)

	fmt.Println("success")
}

type user struct {
	ID   int
	Name string
}

func (u user) Marshal() string {
	return strconv.Itoa(u.ID) + ":" + u.Name
}

func UnmarshalUser(data string) (*user, error) {
	idx := strings.Index(data, ":")
	if idx == -1 ||
		idx == 0 ||
		idx == len(data)-1 ||
		strings.Count(data, ":") != 1 {
		return nil, errors.New("UnmarshalUser: data format incorrect")
	} else {
		u := &user{
			ID:   0,
			Name: "",
		}
		u.ID, _ = strconv.Atoi(data[:idx])
		u.Name = data[idx+1:]
		return u, nil
	}
}

func MustUnmarshalUser(data string) *user {
	u, err := UnmarshalUser(data)
	if err != nil {
		panic(err)
	}
	return u
}

func testUnmarshal(dispatcher *session.Dispatcher) {
	var (
		usr = user{
			ID:   1,
			Name: "tony",
		}.Marshal()
		//gotUsr *user
		got interface{}
	)
	token, err := dispatcher.MarshalToken(usr)
	if err != nil {
		panic(err)
	} else {
		got, _, err = dispatcher.UnmarshalToken(token)
		if err != nil {
			panic(err)
		}
		if usr != got {
			panic(fmt.Errorf("need: %v, got: %v", usr, got))
		}
	}
}

func testExpired(dispatcher *session.Dispatcher) {
	var (
		ID = "expire_id"
	)
	env := dispatcher.GetEnv()
	env.Expired = 5 // 设置过期时间 5 秒
	dispatcher.SetEnv(env)
	token, err := dispatcher.MarshalToken(ID)
	if err != nil {
		panic(err)
	} else {
		// 将时间向前推进 6 秒
		setForwardTime(6*time.Second, func() {
			_, _, err = dispatcher.UnmarshalToken(token)
			if session.ErrInvalidToken != err {
				panic("need invalid error")
			}
		})
	}
}

// 将时间向前推进
func setForwardTime(d time.Duration, done func()) {
	jwt.TimeFunc = func() time.Time {
		return time.Now().Add(d)
	}
	session.NowFunc = func() time.Time {
		return time.Now().Add(d)
	}
	redisMockServer.FastForward(d)
	done()
	jwt.TimeFunc = time.Now
	session.NowFunc = time.Now
}
