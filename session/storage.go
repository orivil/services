// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

import (
	"github.com/orivil/service"
	"time"
)

type Storage interface {

	// 判断 session 是否合法
	IsOK(session string) (ok bool, err error)

	// 创建/更新会话状态
	SaveSession(session string, expires time.Duration) error

	// 移除会话状态
	DelSession(session string) error
}

type StorageService interface {
	service.Provider
	Get(ctn *service.Container) (Storage, error)
}
