// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"fmt"
	"github.com/orivil/service"
	email_captcha "github.com/orivil/services/captcha/email"
	"github.com/orivil/services/captcha/storages/memory"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/email"
	"html/template"
	"io"
	"log"
	"net/http"
)

var config = `
# SMTP 邮箱配置
[smtp-email]

# 域名
host = "smtp.qq.com"

# 端口
port = ":25"

# 用户名
username = "xxx@qq.com"

# 密码
password = "xxx"

# 来源
from = "xxx@qq.com"

# 邮箱验证码配置
[email-captcha]

# 邮件内容类型
content_type = "text/html; charset=UTF-8"

# 验证码过期时间(单位：秒)
expires = 300

# 验证码长度
captcha_length = 4

# 字符集, 最好排除容易混淆的字符. 如果有小写字母，则将被自动转换为大写字母
chars = "2345789ABCDEFHJKLMNPRSTWXYZ"
`

func main() {
	cfgService := cfg.NewService(cfg.NewMemoryStorageService(config))
	emailService := email.NewService("smtp-email", cfgService)
	captchaStorageService := captcha_memory_storage.NewService()
	emailCaptchaService := email_captcha.NewService("email-captcha", cfgService, emailService, captchaStorageService, email_captcha.TemplateMemoryStorage(emailTemplate))
	container := service.NewContainer()
	ds, err := emailCaptchaService.Get(container)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write(indexPage)
	})
	http.HandleFunc("/send", func(writer http.ResponseWriter, request *http.Request) {
		em := request.FormValue("email")
		err := ds.SendCaptcha("verify your email", em)
		if err != nil {
			writer.Write([]byte(err.Error()))
		} else {
			http.Redirect(writer, request, "/verify?email="+em, 302)
		}
	})
	verifyTpl := template.Must(template.New("verify").Parse(verifyPage))
	http.HandleFunc("/verify", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		switch request.Method {
		case http.MethodGet:
			em := request.URL.Query().Get("email")
			verifyTpl.Execute(writer, em)
		case http.MethodPost:
			em := request.FormValue("email")
			code := request.FormValue("code")
			ok, err := ds.Verify(em, code)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			if !ok {
				io.WriteString(writer, "验证码错误!\n")
			} else {
				io.WriteString(writer, "验证通过.\n")
			}
			io.WriteString(writer, "<br><a href='/'>返回</a>")
		}
	})
	fmt.Println("Server is at http://localhost:8666")
	if err = http.ListenAndServe("localhost:8666", nil); err != nil {
		log.Fatal(err)
	}
}

var emailTemplate = `<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Verify Email</title>
</head>
<body>
<h1 style="text-align: center">Verify your Email</h1>
<p><span style="font-weight: bold">{{.}}</span> is your captcha</p>
</body>
</html>`

var indexPage = []byte(`<!doctype html>
<html lang="en">
<head>
    <title>Send Captcha</title>
</head>
<body>
<form action="/send" method=post>
    <p>请输入邮箱:</p>
    <input type="email" name="email">
    <input type="submit" value="获取验证码">
</form>
</body>
</html>`)

var verifyPage = `<!doctype html>
<html lang="en">
<head>
    <title>Verify Captcha</title>
</head>
<body>
<form action="/verify" method=post>
    <p>请输入验证码:</p>
    <input type=hidden name=email value="{{.}}">
    <input type="text" name="code">
    <input type="submit" value="提交">
</form>
</body>
</html>`
