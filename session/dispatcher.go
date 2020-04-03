// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

import (
	"errors"
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

var NowFunc = func() time.Time {
	return time.Now()
}

type WarpError struct {
	original error
	new      error
}

type Claims struct {
	Value interface{}
	jwt.StandardClaims
}

func (w WarpError) Error() string {
	return w.new.Error()
}

func (w WarpError) Unwrap() error {
	return w.original
}

var errInvalidToken = errors.New("invalid token")

func IsInvalidTokenErr(err error) bool {
	if we, ok := err.(WarpError); ok {
		return we.new == errInvalidToken
	} else {
		return err == errInvalidToken
	}
}

type JWTSignMethod string

type Dispatcher struct {
	store Storage
	env   *Env
	mu    sync.Mutex
}

func NewDispatcher(store Storage, env Env) *Dispatcher {
	return &Dispatcher{
		store: store,
		env:   &env,
		mu:    sync.Mutex{},
	}
}

// 生成 token
func (d *Dispatcher) MarshalToken(value interface{}) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.marshalToken(value)
}

func (d *Dispatcher) marshalToken(value interface{}) (string, error) {
	expires := time.Duration(d.env.Expired) * time.Second
	now := NowFunc()
	expiresAt := now.Add(expires).Unix()
	claims := &Claims{
		Value: value,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			NotBefore: now.Unix(),
		},
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
	token, err := ts.SignedString([]byte(d.env.SigningKey))
	if err != nil {
		return "", err
	} else {
		err = d.store.SaveSession(token, expires)
		if err != nil {
			return "", err
		} else {
			return token, nil
		}
	}
}

// 解析 ID, 如果 token 解析失败，可使用 IsInvalidTokenErr 判断 token 是否需要重新生成.
// 返回值 value 为生成 token 时传入的数据, expireAt 为 token 过期时间
func (d *Dispatcher) UnmarshalToken(token string) (value interface{}, expireAt int64, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if token == "" {
		return nil, 0, errInvalidToken
	}
	var ok bool
	ok, err = d.store.IsOK(token)
	if err != nil {
		return nil, 0, err
	}
	if !ok {
		return nil, 0, errInvalidToken
	}

	var jt *jwt.Token
	jt, err = jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(d.env.SigningKey), nil
	})
	if err != nil {
		return nil, 0, WarpError{
			original: err,
			new:      errInvalidToken,
		}
	} else if !jt.Valid {
		return nil, 0, errInvalidToken
	}
	claims := jt.Claims.(*Claims)
	return claims.Value, claims.ExpiresAt, nil
}

func (d *Dispatcher) GetEnv() *Env {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.env
}

func (d *Dispatcher) SetEnv(env *Env) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.env = env
}

func (d *Dispatcher) DelToken(token string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.store.DelSession(token)
}