// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"fmt"
	"github.com/orivil/log"
	"github.com/orivil/service"
	"github.com/orivil/services/auth/oauth2/wechat"
	wechat_storages "github.com/orivil/services/auth/oauth2/wechat/storages"
	"github.com/orivil/services/cfg"
	"net/http"
)

var toml = `# 公众号登录授权
[oauth2-wechat]
# 公众号 AppID
appid = "wx4d2b9900882b9fe6"
# 公众号 App Secret
secret = "8b70dabf6cfc479b4fb497d454c95173"
# 授权成功后跳转地址, 在发起授权时可以指定地址, 如果未指定才使用该值
redirect_uri = ""
`

func main() {
	container := service.NewContainer()
	configService := cfg.NewService(cfg.NewMemoryStorageService(toml))
	wechatService := wechat.NewService(wechat.Options{
		ConfigNamespace: "oauth2-wechat",
		ConfigService:   configService,
		Storage:         wechat_storages.NewMemoryService(),
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
		log.Info.Println(uri)
		http.Redirect(writer, request, uri, 302)
	})
	http.HandleFunc("/signing", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		code := query.Get("code")
		if code != "" {
			openid, err := dis.Exchange(code)
			if err != nil {
				panic(err)
			}
			user, err := dis.GetUserInfo(openid)
			if err != nil {
				panic(err)
			}
			fmt.Println(user)
		}
	})
	addr := ":8080"
	log.Info.Printf("listen and serve http://localhost%s\n", addr)
	err = http.ListenAndServe(addr, nil)
	log.Panic.Println(err)
}
