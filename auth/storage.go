// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package auth

type Storage interface {
	// 获得账号已加密密码, 如果账号不存在则返回空密码
	GetPassword(usernameOrEmail string) (password string, err error)

	// 保存账号密码
	SavePassword(username, email, password string) error
}
