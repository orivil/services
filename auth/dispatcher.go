// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package auth

import (
	"errors"
	"github.com/orivil/services/session"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

var (
	ErrUserNotRegistered     = errors.New("user not registered")
	ErrUserAlreadyRegistered = errors.New("user already registered")
	ErrPasswordIncorrect     = errors.New("password incorrect")
	ErrInvalidToken          = session.ErrInvalidToken
)

type Dispatcher struct {
	store Storage
	jwt   *session.JWTAuth
}

func NewDispatcher(store Storage, jwt *session.JWTAuth) *Dispatcher {
	return &Dispatcher{
		store: store,
		jwt:   jwt,
	}
}

// 更新账号密码
func (h *Dispatcher) UpdatePassword(id int, password string) error {
	exist, err := h.store.GetPassword(id)
	if err != nil {
		return err
	}
	if exist != "" {
		return h.generateAndSavePassword(id, password)
	} else {
		return ErrUserNotRegistered
	}
}

func (h *Dispatcher) generateAndSavePassword(id int, password string) error {
	hashedPass, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err1 != nil {
		return err1
	}
	return h.store.SavePassword(id, string(hashedPass))
}

// 注册账号, 如果账号已存在则返回 ErrUserAlreadyRegistered 错误
func (h *Dispatcher) Register(id int, password string) error {
	exist, err := h.store.GetPassword(id)
	if err != nil {
		return err
	}
	if exist != "" {
		return ErrUserAlreadyRegistered
	}
	return h.generateAndSavePassword(id, password)
}

// 登录账号并设置登录状态. 如果用户未注册则返回 ErrUserNotRegistered 错误,
// 密码错误则返回 ErrPasswordIncorrect
func (h *Dispatcher) Login(id int, password string) (token string, err error) {
	hashedPassword, er := h.store.GetPassword(id)
	if er != nil {
		return "", er
	}
	if hashedPassword == "" {
		return "", ErrUserNotRegistered
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", ErrPasswordIncorrect
	}
	return h.jwt.MarshalToken(strconv.Itoa(id))
}

// 解除登录状态
func (h *Dispatcher) Logout(token string) error {
	return h.jwt.DelToken(token)
}

// 解析 token 获得用户 ID, 如果获得了新的 newToken, 则表示原 token 快要过期, 需要替换 token
// 如果没有获得 newToken, 则原 token 可以继续使用. 如果 token 验证错误, 则返回 ErrInvalidToken
// 此时需要重新登录获得 token
func (h *Dispatcher) GetUserID(token string) (id int, newToken string, err error) {
	var idStr string
	idStr, newToken, err = h.jwt.UnmarshalID(token)
	if err == session.ErrInvalidToken {
		err = ErrInvalidToken
	}
	id, err = strconv.Atoi(idStr)
	return
}
