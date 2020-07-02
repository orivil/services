// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package wechat_storages

import (
	"github.com/orivil/services/auth/oauth2/wechat"
	"github.com/orivil/wechat/oauth2"
)

type Memory struct {
	users map[string]*oauth2.User
	tokens map[string]*wechat.Token
}

func NewMemory() *Memory {
	return &Memory{
		users: map[string]*oauth2.User{},
		tokens: map[string]*wechat.Token{},
	}
}

func (m *Memory) SaveUser(openid string, user *oauth2.User) error {
	m.users[openid] = user
	return nil
}

func (m *Memory) GetUser(openid string) (*oauth2.User, error) {
	user := m.users[openid]
	return user, nil
}

func (m *Memory) SaveToken(openid string, token *wechat.Token) error {
	m.tokens[openid] = token
	return nil
}

func (m *Memory) GetToken(openid string) (*wechat.Token, error) {
	token := m.tokens[openid]
	return token, nil
}

