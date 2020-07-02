// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package wechat

import (
	"github.com/orivil/service"
	"github.com/orivil/wechat/oauth2"
)

type Storage interface {
	// store token
	SaveToken(openid string, token *Token) error
	GetToken(openid string) (*Token, error)

	// store user
	SaveUser(openid string, user *oauth2.User) error
	GetUser(openid string) (*oauth2.User, error)
}

type StorageService interface {
	service.Provider
	Get(ctn *service.Container) (Storage, error)
}