// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

import "time"

type Storage interface {
	// 获得在线会话人数(不同地方登录在线数)
	GetOnlineSessions(id string) (total int, err error)

	// 判断 session 是否合法
	IsOK(id, session string) (ok bool, err error)

	// 创建/更新会话状态
	SaveSession(id, session string, expires time.Duration) error

	// 移除会话状态
	DelSession(id, session string) error

	// 移除最先过期的会话状态
	DelFirstExpireSession(id string) error
}
