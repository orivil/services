// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

// build ignore

package main

import (
	"fmt"
	"github.com/orivil/service"
	image_captcha "github.com/orivil/services/captcha/image"
	"github.com/orivil/services/captcha/storages/memory"
	"github.com/orivil/services/cfg"
	"io"
	"log"
	"net/http"
	"path"
	"text/template"
)

var config = `
# 验证码配置
[image-captcha]

# 验证码图片宽度(单位: px)
img_width = 240

# 验证码图片高度(单位：px)
img_height = 80

# 验证码长度
captcha_length = 4

# 验证码过期时间(单位：秒)
expires = 300

# 字符集, 最好排除容易混淆的字符. 如果有小写字母，则将被自动转换为大写字母
chars = "2345789ABCDEFHJKLMNPRSTWXYZ"
`

var formTemplate = template.Must(template.New("example").Parse(formTemplateSrc))

func main() {
	cfgService := cfg.NewService(cfg.NewMemoryStorageService(config))
	captchaStorageService := captcha_memory_storage.NewService()
	captchaDispatcherService := image_captcha.NewService("image-captcha", cfgService, captchaStorageService)

	container := service.NewContainer()

	cd, err := captchaDispatcherService.Get(container)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			http.NotFound(writer, request)
			return
		}
		d := struct {
			CaptchaId string
		}{
			cd.GenID(),
		}
		if err = formTemplate.Execute(writer, &d); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/verify", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		ok, err := cd.Verify(request.FormValue("captchaId"), request.FormValue("captchaSolution"))
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
	})
	http.HandleFunc("/captcha/", func(writer http.ResponseWriter, request *http.Request) {
		_, file := path.Split(request.URL.Path)
		ext := path.Ext(file)
		id := file[:len(file)-len(ext)]
		if ext == "" || id == "" {
			http.NotFound(writer, request)
			return
		}
		err = cd.WriteCodeToImage(id, writer)
		if err != nil {
			log.Println(err)
		}
	})
	fmt.Println("Server is at http://localhost:8666")
	if err = http.ListenAndServe("localhost:8666", nil); err != nil {
		log.Fatal(err)
	}
}

var formTemplateSrc = `<!doctype html>
<html lang="en">
<head>
    <title>Captcha Test</title>
    <script>
        function setSrcQuery(e, q) {
            let src  = e.src;
            let p = src.indexOf('?');
            if (p >= 0) {
                src = src.substr(0, p);
            }
            e.src = src + "?" + q
        }
        function reload() {
            setSrcQuery(document.getElementById('image'), "reload=" + (new Date()).getTime());
            return false;
        }
    </script>
</head>
<body>
<form action="/verify" method=post>
    <p>请输入验证码:</p>
    <p>
        <a href="javascript:void(0)" onclick="reload()">
            <img id=image src="/captcha/{{.CaptchaId}}.png" alt="Captcha image">
        </a>
    </p>
    <input type=hidden name=captchaId value="{{.CaptchaId}}"><br>
    <input name=captchaSolution>
    <input type=submit value=提交>
</form>
</body>
</html>`
