// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session_memory_storage

import (
	"fmt"
	"github.com/orivil/service"
	"time"
)

func ExampleNewService() {
	storeService := NewService()
	container := service.NewContainer()
	store, _ := storeService.Get(container)
	session := "tokenStr"

	err := store.SaveSession(session, 2*time.Second)
	panicErr(err)
	var ok bool
	ok, err = store.IsOK(session)
	panicErr(err)
	fmt.Println(ok == true)

	// 将时间向前推进, 以删除过期数据
	NowFunc = func() time.Time {
		return time.Now().Add(4 * time.Second)
	}

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
