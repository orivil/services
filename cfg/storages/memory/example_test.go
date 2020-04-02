// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package memory_test

import (
	"fmt"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/services/cfg/storages/memory"
)

// 内存型配置存储器主要用于测试时方做便零时配置
func ExampleNewService() {

	var configData = `
[mysql]
host = "127.0.0.1"
port= "3306"
`
	configService := cfg.NewService(memory.NewService(configData))

	container := service.NewContainer()

	envs, err := configService.Get(container)
	if err != nil {
		panic(err)
	}
	envs, err = envs.GetSub("mysql")
	for s, i := range envs {
		fmt.Printf("%s = %s\n", s, i)
	}
	// Output:
	// host = 127.0.0.1
	// port = 3306
}
