// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package wechat

import (
	"errors"
	"github.com/orivil/wechat/oauth2"
	"sync"
	"time"
)

const refreshExpires = 29 * 24 * 60 * 60 // 29 天之内允许刷新 token (实际时间是30天)

var ErrTokenExpired = errors.New("token expired")

type Dispatcher struct {
	cfg *Config
	store Storage
	mu sync.Mutex
}

func NewDispatcher(cfg *Config, store Storage) *Dispatcher {
	return &Dispatcher{
		cfg: cfg,
		store: store,
		mu:  sync.Mutex{},
	}
}

// 获得授权地址, 用户授权后会跳转至地址 url, 如果不提供 url 参数, 则将配置项 redirect_uri 作为跳转地
// 址, 并带上 code, state 参数
func (d *Dispatcher) RedirectURI(scope Scope, state, url string) string {
	return d.cfg.RedirectURL(string(scope), state, url)
}

// 使用 code 换取 access token, 并保存 access token, 返回用户 openid
func (d *Dispatcher) Exchange(code string) (token *Token, err error) {
	token, err = d.cfg.Exchange(code)
	if err != nil {
		return
	}
	err = d.store.SaveToken(token.Openid, token)
	if err != nil {
		return
	}
	return token, nil
}

// 获取 access token, 且尝试刷新 access token
func (d *Dispatcher) GetToken(openid string) (token *Token, err error) {
	token, err = d.store.GetToken(openid)
	if err != nil {
		return nil, err
	}
	if token != nil {
		if now := time.Now().Unix(); token.ExpiresAt < now {
			if token.ExpiresAt + refreshExpires < now {
				return nil, ErrTokenExpired
			} else {
				// 刷新 token
				token, err = d.cfg.Refresh(token.RefreshToken)
				if err != nil {
					return nil, err
				}
				err = d.store.SaveToken(openid, token)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		return nil, ErrTokenExpired
	}
	return token, nil
}

// 获得用户信息, 需用户允许(配置项 scope = snsapi_userinfo 且用户同意), 如果 access token
// 已过期, 则会刷新 access token, 如果 access token 超过刷新时间, 则会返回 ErrTokenExpired
func (d *Dispatcher) GetUserInfo(openid string) (user *oauth2.User, err error) {
	user, err = d.store.GetUser(openid)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	var token *Token
	token, err = d.GetToken(openid)
	if err != nil {
		return nil, err
	}
	return d.cfg.GetUserInfo(openid, token.AccessToken)
}