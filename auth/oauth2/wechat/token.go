// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package wechat

import (
	"github.com/orivil/wechat/oauth2"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

func NewToken(t *oauth2.AccessToken) *Token {
	return &Token {
		AccessToken:  t.AccessToken,
		ExpiresAt:    time.Now().Unix() + t.ExpiresIn,
		RefreshToken: t.RefreshToken,
		Openid:       t.Openid,
		Scope:        t.Scope,
	}
}