// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
)

const (
	JWTSignMethodHS512 JWTSignMethod = "HS512"
	JWTSignMethodHS384 JWTSignMethod = "HS384"
	JWTSignMethodHS256 JWTSignMethod = "HS256"
)

var ErrInvalidToken = errors.New("invalid token")

type JWTSignMethod string

/**
# 用户认证(jwt)
[jwt_auth]
# 同一账号同时在线客户端数量, 如果小于 0 则不限制
max_online_sessions = 3
# 签名方式(支持：HS512/HS384/HS256)
signing_method = "HS512"
# 签名key
signing_key = "GINADMIN"
# 过期时间（单位秒）
expired = 7200
# 存储(支持：file/redis)
store = "file"
# 文件路径
file_path = "data/jwt_auth.db"
# redis数据库(如果存储方式是redis，则指定存储的数据库)
redis_db = 10
# 存储到redis数据库中的键名前缀
redis_prefix = "auth_"
*/

type JWTEnv struct {
	MaxOnlineSessions int    `toml:"max_online_sessions"`
	SigningMethod     string `toml:"signing_method"`
	SigningKey        string `toml:"signing_key"`
	Expired           int    `toml:"expired"`
	Store             string `toml:"store"`
	FilePath          string `toml:"file_path"`
	RedisDB           string `toml:"redis_db"`
	RedisPrefix       string `toml:"redis_prefix"`
}

type JWTAuth struct {
	store Storage
	env   *JWTEnv
	mu    sync.Mutex
}

func (j *JWTAuth) MarshalToken(id string) (string, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	expires := time.Duration(j.env.Expired) * time.Second
	now := time.Now()
	expiresAt := now.Add(expires).Unix()
	claims := jwt.StandardClaims{
		Audience:  "",
		ExpiresAt: expiresAt,
		Id:        "",
		IssuedAt:  0,
		Issuer:    "",
		NotBefore: now.Unix(),
		Subject:   id,
	}
	var sm jwt.SigningMethod
	switch JWTSignMethod(j.env.SigningMethod) {
	case JWTSignMethodHS512:
		sm = jwt.SigningMethodHS512
	case JWTSignMethodHS384:
		sm = jwt.SigningMethodHS384
	case JWTSignMethodHS256:
		sm = jwt.SigningMethodHS256
	}
	ts := jwt.NewWithClaims(sm, claims)
	token, err := ts.SignedString(j.env.SigningKey)
	if err != nil {
		return "", err
	} else {
		if j.env.MaxOnlineSessions > 0 {
			var total, err1 = j.store.GetOnlineSessions(id)
			if err1 != nil {
				return "", err1
			}
			if total >= j.env.MaxOnlineSessions {
				err = j.store.DelFirstExpireSession(id)
				if err != nil {
					return "", err
				}
			}
		}
		err = j.store.SaveSession(id, token, expires)
		if err != nil {
			return "", err
		} else {
			return token, nil
		}
	}
}

// 解析 ID, 如果 token 解析失败，返回 ErrInvalidToken 错误, expireAt 为 token 过期时间，
// 可以在 token 过期之前重新生成 token
func (j *JWTAuth) UnmarshalID(token string) (id string, expireAt int64, err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if token == "" {
		return "", 0, ErrInvalidToken
	}
	var t, err1 = jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(j.env.SigningKey), nil
	})
	if err1 != nil {
		return "", 0, err1
	} else if !t.Valid {
		return "", 0, ErrInvalidToken
	}
	claims := t.Claims.(*jwt.StandardClaims)
	id = claims.Subject
	var ok, er = j.store.IsOK(id, token)
	if er != nil {
		return "", 0, err
	}
	if !ok {
		return "", 0, ErrInvalidToken
	} else {
		return id, claims.ExpiresAt, nil
	}
}

func (j *JWTAuth) RemoveToken(id, token string) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.store.DelSession(id, token)
}
