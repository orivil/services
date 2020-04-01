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

type JWTAuth struct {
	store Storage
	env   *JWTEnv
	mu    sync.Mutex
}

func NewJWTAuth(store Storage, env JWTEnv) *JWTAuth {
	return &JWTAuth{
		store: store,
		env:   &env,
		mu:    sync.Mutex{},
	}
}

// 生成 token
func (j *JWTAuth) MarshalToken(id string) (string, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.marshalToken(id)
}

func (j *JWTAuth) marshalToken(id string) (string, error) {
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
		err = j.store.SaveSession(token, expires)
		if err != nil {
			return "", err
		} else {
			return token, nil
		}
	}
}

// 解析 ID, 如果 token 解析失败，返回 ErrInvalidToken 错误, 返回值 newToken 为重新生成的 token,
// 在过期之前的一段时间自动生成
func (j *JWTAuth) UnmarshalID(token string) (id, newToken string, err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if token == "" {
		return "", "", ErrInvalidToken
	}
	var t, err1 = jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(j.env.SigningKey), nil
	})
	if err1 != nil {
		return "", "", err1
	} else if !t.Valid {
		return "", "", ErrInvalidToken
	}
	claims := t.Claims.(*jwt.StandardClaims)
	id = claims.Subject
	var ok, er = j.store.IsOK(token)
	if er != nil {
		return "", "", err
	}
	if !ok {
		return "", "", ErrInvalidToken
	} else {
		newToken, err = j.refresh(id, claims.ExpiresAt)
		return id, newToken, err
	}
}

// 判断 token 是否需要重新生成，如果需要则生成新的 token
func (j *JWTAuth) refresh(id string, tokenExpireAt int64) (token string, err error) {
	if tokenExpireAt < time.Now().Add(time.Duration(j.env.Refresh)*time.Second).Unix() {
		return j.marshalToken(id)
	} else {
		return "", nil
	}
}

func (j *JWTAuth) DelToken(token string) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.store.DelSession(token)
}
