// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package wechat_storages

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/orivil/wechat/oauth2"
	"time"
)

var NilContext = context.TODO()

var tokenDuration = time.Hour * 24 * 30 // token 保存 30 天

type Redis struct {
	tc *redis.Client
	uc *redis.Client
}

func NewRedis(token, user *redis.Client) *Redis {
	return &Redis{
		tc: token,
		uc: user,
	}
}

func (r *Redis) SaveToken(openid string, token *oauth2.AccessToken) error {
	str, err := jsonStr(token)
	if err != nil {
		return err
	}
	return r.tc.Set(NilContext, openid, str, tokenDuration).Err()
}

func (r *Redis) GetToken(openid string) (*oauth2.AccessToken, error) {
	str, err := r.tc.Get(NilContext, openid).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	token := &oauth2.AccessToken{}
	err = json.Unmarshal([]byte(str), token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *Redis) SaveUser(openid string, user *oauth2.User) error {
	str, err := jsonStr(user)
	if err != nil {
		return err
	}
	return r.uc.Set(NilContext, openid, str, tokenDuration).Err()
}

func (r *Redis) GetUser(openid string) (*oauth2.User, error) {
	str, err := r.uc.Get(NilContext, openid).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	user := &oauth2.User{}
	err = json.Unmarshal([]byte(str), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func jsonStr(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}