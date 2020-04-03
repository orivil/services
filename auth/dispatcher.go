// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package auth

import (
	"errors"
	"github.com/orivil/services/captcha"
	"github.com/orivil/services/session"
	"golang.org/x/crypto/bcrypt"
	"io"
)

var (
	ErrUserNotRegistered     = errors.New("user not registered")
	ErrUserAlreadyRegistered = errors.New("user already registered")
	ErrPasswordIncorrect     = errors.New("password incorrect")
	ErrInvalidToken          = errors.New("invalid token")
)

type Dispatcher struct {
	store             Storage
	sessionDispatcher *session.Dispatcher
	captchaDispatcher *captcha.Dispatcher
}

func NewDispatcher(store Storage, sessionDispatcher *session.Dispatcher, captchaDispatcher *captcha.Dispatcher) *Dispatcher {
	return &Dispatcher{
		store:             store,
		sessionDispatcher: sessionDispatcher,
		captchaDispatcher: captchaDispatcher,
	}
}

func (h *Dispatcher) GetCaptchaID() (id string, err error) {
	return h.captchaDispatcher.GenID()
}

func (h *Dispatcher) ServeImage(id string, w io.Writer) error {
	return h.captchaDispatcher.WriteImage(id, w)
}

func (h *Dispatcher) VerifyCaptcha(id, code string) (ok bool, err error) {
	return h.captchaDispatcher.Verify(id, code)
}

// 更新账号密码, 用户不存在则返回 ErrUserNotRegistered 错误
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
func (h *Dispatcher) CreatePassword(id int, password string) error {
	exist, err := h.store.GetPassword(id)
	if err != nil {
		return err
	}
	if exist != "" {
		return ErrUserAlreadyRegistered
	}
	return h.generateAndSavePassword(id, password)
}

// 检测密码是否正确. 如果用户未注册则返回 ErrUserNotRegistered 错误,
// 密码错误则返回 ErrPasswordIncorrect
func (h *Dispatcher) VerifyPassword(id int, password string) (err error) {
	var hashedPassword string
	hashedPassword, err = h.store.GetPassword(id)
	if err != nil {
		return err
	}
	if hashedPassword == "" {
		return ErrUserNotRegistered
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrPasswordIncorrect
	} else {
		return nil
	}
}

// 生成 token 并设置会话, 在操作之前应该验证用户密码及验证码等
func (h *Dispatcher) MarshalToken(user interface{}) (token string, err error) {
	return h.sessionDispatcher.MarshalToken(user)
}

// 删除 token 并移除会话
func (h *Dispatcher) DelToken(token string) error {
	return h.sessionDispatcher.DelToken(token)
}

// 解析 token 并获得 token 过期时间, 如果返回值 err 为 ErrInvalidToken, 则需要重新验证并生成 token
func (h *Dispatcher) UnmarshalToken(token string) (user interface{}, expiredAt int64, err error) {
	user, expiredAt, err = h.sessionDispatcher.UnmarshalToken(token)
	if session.IsInvalidTokenErr(err) {
		err = ErrInvalidToken
	}
	return
}
