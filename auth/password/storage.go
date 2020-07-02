// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package password

import "github.com/orivil/service"

type Storage interface {
	// 获得账号已加密密码, 如果账号不存在则返回空密码
	GetPassword(id int) (password string, err error)

	// 保存账号密码
	SavePassword(id int, password string) error
}

type StorageService interface {
	service.Provider
	Get(ctn *service.Container) (Storage, error)
}
