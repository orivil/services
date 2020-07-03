// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
)

const (
	JWTSignMethodHS512 JWTSignMethod = "HS512"
	JWTSignMethodHS384 JWTSignMethod = "HS384"
	JWTSignMethodHS256 JWTSignMethod = "HS256"
)

var NowFunc = time.Now

type JWTSignMethod string

type Dispatcher struct {
	store Storage
	env   *Env
	mu    sync.Mutex
}

func NewDispatcher(store Storage, env Env) *Dispatcher {
	return &Dispatcher {
		store: store,
		env:   &env,
		mu:    sync.Mutex{},
	}
}

func (d *Dispatcher) MarshalToken(id string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := NowFunc().Unix()
	expireAt := d.env.Expires + now
	claims := &jwt.StandardClaims {
		ExpiresAt: expireAt,
		NotBefore: now,
		Id: id,
	}
	var sm jwt.SigningMethod
	switch JWTSignMethod(d.env.SigningMethod) {
	case JWTSignMethodHS512:
		sm = jwt.SigningMethodHS512
	case JWTSignMethodHS384:
		sm = jwt.SigningMethodHS384
	case JWTSignMethodHS256:
		sm = jwt.SigningMethodHS256
	default:
		return "", fmt.Errorf("jwt_auth config field: signing_method = %s is not allowed", d.env.SigningMethod)
	}
	ts := jwt.NewWithClaims(sm, claims)
	jwtToken, err := ts.SignedString([]byte(d.env.SigningKey))
	if err != nil {
		return "", err
	} else {
		err = d.store.SaveSession(jwtToken, time.Duration(d.env.Expires) * time.Second)
		if err != nil {
			return "", err
		} else {
			return jwtToken, nil
		}
	}
}

// 解析 jwt token, 如果解释失败(包括过期), 则返回 nil token, 发生错误则返回 error
func (d *Dispatcher) UnmarshalToken(jwtToken string) (token *Token, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if jwtToken == "" {
		return nil, nil
	}
	var ok bool
	ok, err = d.store.IsOK(jwtToken)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	var jt *jwt.Token
	jt, err = jwt.ParseWithClaims(jwtToken, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(d.env.SigningKey), nil
	})
	if jt == nil || !jt.Valid {
		return nil, nil
	}
	claims := jt.Claims.(*jwt.StandardClaims)
	return &Token{ ID: claims.Id, ExpiresAt: claims.ExpiresAt }, nil
}

func (d *Dispatcher) DelToken(token string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.store.DelSession(token)
}
