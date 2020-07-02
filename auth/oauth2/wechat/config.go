// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package wechat

import "github.com/orivil/wechat/oauth2"

/**
# 公众号登录授权
[oauth2-wechat]
# 公众号 AppID
appid = ""
# 公众号 App Secret
secret = ""
# 授权成功后跳转地址, 在发起授权时可以指定地址, 如果未指定才使用该值
redirect_uri = ""
*/
type Config struct {
	Appid string `toml:"appid"`
	Secret string `toml:"secret"`
	RedirectURI string `toml:"redirect_uri"`
}

func (c *Config) RedirectURL(scope, state, url string) string {
	if url == "" {
		url = c.RedirectURI
	}
	return oauth2.InitAppRedirect(scope, url, c.Appid, "", state)
}

func (c *Config) Exchange(code string) (*oauth2.AccessToken, error) {
	return oauth2.GetAccessToken(c.Appid, c.Secret, code)
}

func (c *Config) Refresh(refreshToken string) (*oauth2.AccessToken, error) {
	return oauth2.RefreshAccessToken(c.Appid, refreshToken)
}

func (c *Config) GetUserInfo(openid, token string) (*oauth2.User, error) {
	return oauth2.GetUserInfo(openid, token)
}
