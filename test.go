// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/orivil/log"
	"github.com/orivil/service"
	"github.com/orivil/services/auth/oauth2/wechat"
	wechat_storages "github.com/orivil/services/auth/oauth2/wechat/storages"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/memory/redis"
	"io/ioutil"
	"net/http"
)

var toml =
`# 公众号登录授权
[oauth2-wechat]
# 公众号 AppID
appid = "wx4d2b9900882b9fe6"
# 公众号 App Secret
secret = "8b70dabf6cfc479b4fb497d454c95173"
# 授权成功后跳转地址, 在发起授权时可以指定地址, 如果未指定才使用该值
redirect_uri = ""

# oauth2 redis user 配置
[oauth2-redis-user]
# 是否模拟客户端, 开启后将链接至虚拟客户端, 测试时使用
# 使用方式参考: https://github.com/alicebob/miniredis
is_mock = false
# 地址
addr = "127.0.0.1:6379"
# 密码
password = ""
# 数据库
db = 1

# oauth2 redis token 配置
[oauth2-redis-token]
# 是否模拟客户端, 开启后将链接至虚拟客户端, 测试时使用
# 使用方式参考: https://github.com/alicebob/miniredis
is_mock = false
# 地址
addr = "127.0.0.1:6379"
# 密码
password = ""
# 数据库
db = 2
`

func main() {
	container := service.NewContainer()
	configService := cfg.NewService(cfg.NewMemoryStorageService(toml))
	oauth2RedisUserService := redis.NewService("oauth2-redis-user", configService)
	oauth2RedisTokenService := redis.NewService("oauth2-redis-token", configService)
	wechatService := wechat.NewService(wechat.Options{
		ConfigNamespace: "oauth2-wechat",
		ConfigService:   configService,
		Storage:         wechat_storages.NewRedisService(oauth2RedisUserService, oauth2RedisTokenService),
	})
	dis, err := wechatService.Get(container)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/wechat/notify", func(writer http.ResponseWriter, request *http.Request) {
		echoStr := request.URL.Query().Get("echostr")
		writer.Write([]byte(echoStr))
	})
	http.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		uri := dis.RedirectURI(wechat.ScopeUserInfo, "some_state", "http://data.orivil.com/signing")
		http.Redirect(writer, request, uri, 302)
	})
	http.HandleFunc("/refresh", func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("RefreshToken")
		if err != nil {
			panic(err)
		}
		refreshToken := cookie.Value
		uri := "https://api.weixin.qq.com/sns/oauth2/access_token?grant_type=refresh_token&refresh_token" + refreshToken
		uri += "&appid=wx4d2b9900882b9fe6&secret=8b70dabf6cfc479b4fb497d454c95173"
		fmt.Println(uri)
		res, err := http.Get(uri)
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	})

	http.HandleFunc("/signing", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		code := query.Get("code")
		if code != "" {
			token, err := dis.Exchange(code)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%+v\n", token)
			user, err := dis.GetUserInfo(token.Openid)
			if err != nil {
				panic(err)
			}
			data, err := json.MarshalIndent(user, "", "	")
			if err != nil {
				panic(err)
			}
			http.SetCookie(writer, &http.Cookie{
				Name: "Authorization",
				Value: token.AccessToken,
			})
			http.SetCookie(writer, &http.Cookie{
				Name: "RefreshToken",
				Value: token.RefreshToken,
			})
			writer.Write(data)
			fmt.Println(user)
		}
	})
	addr := ":8080"
	log.Info.Printf("listen and serve http://localhost%s\n", addr)
	err = http.ListenAndServe(addr, nil)
	log.Panic.Println(err)
}
