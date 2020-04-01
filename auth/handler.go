// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package auth

import (
	"errors"
	"github.com/orivil/services/session"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotRegistered     = errors.New("user not registered")
	ErrUserAlreadyRegistered = errors.New("user already registered")
	ErrPasswordIncorrect     = errors.New("password incorrect")
	ErrInvalidToken          = session.ErrInvalidToken
)

type Handler struct {
	store Storage
	jwt   *session.JWTAuth
}

func NewHandler(store Storage, jwt *session.JWTAuth) *Handler {
	return &Handler{
		store: store,
		jwt:   jwt,
	}
}

// 更新账号密码
func (h *Handler) UpdatePassword(username, email, password string) error {
	exist, err := h.store.GetPassword(username, email)
	if err != nil {
		return err
	}
	if exist != "" {
		return h.generateAndSavePassword(username, email, password)
	} else {
		return ErrUserNotRegistered
	}
}

func (h *Handler) generateAndSavePassword(username, email, password string) error {
	hashedPass, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err1 != nil {
		return err1
	}
	return h.store.SavePassword(username, email, string(hashedPass))
}

// 注册账号, 如果账号已存在则返回 ErrUserAlreadyRegistered 错误
func (h *Handler) Register(account, password string) error {
	exist, err := h.store.GetPassword(account)
	if err != nil {
		return err
	}
	if exist != "" {
		return ErrUserAlreadyRegistered
	}
	return h.generateAndSavePassword(account, password)
}

// 登录账号并设置登录状态, 参数 userID 用于生成 token, 可使用生成的 token 获得 userID
// 如果用户未注册则返回 ErrUserNotRegistered 错误, 密码错误则返回 ErrPasswordIncorrect
func (h *Handler) Login(account, password, userID string) (token string, err error) {
	hashedPassword, er := h.store.GetPassword(account)
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
	return h.jwt.MarshalToken(userID)
}

// 解除登录状态
func (h *Handler) Logout(token string) error {
	return h.jwt.DelToken(token)
}

// 解析 token 获得用户 ID, 如果获得了新的 newToken, 则表示原 token 快要过期, 需要替换 token
// 如果没有获得 newToken, 则原 token 可以继续使用. 如果 token 验证错误, 则返回 ErrInvalidToken
// 此时需要重新登录获得 token
func (h *Handler) GetUserID(token string) (id, newToken string, err error) {
	id, newToken, err = h.jwt.UnmarshalID(token)
	if err == session.ErrInvalidToken {
		err = ErrInvalidToken
	}
	return
}
