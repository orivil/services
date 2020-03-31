// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg_test

import (
	"fmt"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
	"reflect"
)

type Mysql struct {
	Host       string `toml:"host"`
	Port       string `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	DBName     string `toml:"db_name"`
	Parameters string `toml:"parameters"`
}

func ExampleProvider() {
	// 可重新设置配置文件地址
	cfg.ConfigFile = "configs/config.toml"

	// 依赖容器
	container := service.NewContainer(false)

	// 从容器中获得配置数据，该容器中的配置只会初始化一次
	env := container.MustGet(&cfg.Service).(xcfg.Env)

	// 解析配置数据
	mysql := &Mysql{}
	err := env.UnmarshalSub("mysql", mysql)
	if err != nil {
		panic(err)
	}

	var needMysql = Mysql{
		Host:       "127.0.0.1",
		Port:       "3306",
		User:       "root",
		Password:   "123456",
		DBName:     "ginadmin",
		Parameters: "charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
	}

	// 测试结果
	fmt.Println(reflect.DeepEqual(*mysql, needMysql))

	// Output:
	// true
}
