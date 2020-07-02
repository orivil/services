// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package password

import (
	"errors"
	"github.com/orivil/limiter"
	"github.com/orivil/services/captcha/email"
	"github.com/orivil/services/captcha/image"
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
	Store                  Storage
	SessionDispatcher      *session.Dispatcher
	ImageCaptchaDispatcher *image_captcha.Dispatcher
	EmailCaptchaDispatcher *email_captcha.Dispatcher
	Limiter                *limiter.TimesLimiter
}

// 发送邮箱验证码
func (h *Dispatcher) SendEmailCaptcha(subject, email string) error {
	return h.EmailCaptchaDispatcher.SendCaptcha(subject, email)
}

// 验证邮箱验证码
func (h *Dispatcher) VerifyEmailCaptcha(email, code string) (ok bool, err error) {
	return h.EmailCaptchaDispatcher.Verify(email, code)
}

// 获得图片验证码 ID
func (h *Dispatcher) GetImageCaptchaID() string {
	return h.ImageCaptchaDispatcher.GenID()
}

// 写入图片验证码到接口中
func (h *Dispatcher) ServeCaptchaImage(id string, w io.Writer) error {
	return h.ImageCaptchaDispatcher.WriteCodeToImage(id, w)
}

// 验证图片验证码
func (h *Dispatcher) VerifyImageCaptcha(id, code string) (ok bool, err error) {
	return h.ImageCaptchaDispatcher.Verify(id, code)
}

// 注册账号, 如果账号已存在则返回 ErrUserAlreadyRegistered 错误
func (h *Dispatcher) CreatePassword(id int, password string) error {
	exist, err := h.Store.GetPassword(id)
	if err != nil {
		return err
	}
	if exist != "" {
		return ErrUserAlreadyRegistered
	}
	return h.generateAndSavePassword(id, password)
}

// 更新账号密码, 用户不存在则返回 ErrUserNotRegistered 错误
func (h *Dispatcher) UpdatePassword(id int, password string) error {
	exist, err := h.Store.GetPassword(id)
	if err != nil {
		return err
	}
	if exist != "" {
		return h.generateAndSavePassword(id, password)
	} else {
		return ErrUserNotRegistered
	}
}

// 检测密码是否正确. 如果用户未注册则返回 ErrUserNotRegistered 错误,
// 密码错误则返回 ErrPasswordIncorrect
func (h *Dispatcher) VerifyPassword(id int, password string) (err error) {
	var hashedPassword string
	hashedPassword, err = h.Store.GetPassword(id)
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

// 生成 token 并设置会话, 参数 user 可以是用户 ID, 或序列化后的字符串, 参数 value 最好不要太复杂,
// 避免影响性能(包括本地 session 合法性验证, 前后端数据传输过程中消耗的性能). 在操作之前应该验证用
// 户密码及验证码等
func (h *Dispatcher) MarshalToken(user string) (token string, err error) {
	return h.SessionDispatcher.MarshalToken(user)
}

// 删除 token 并移除会话
func (h *Dispatcher) DelToken(token string) error {
	return h.SessionDispatcher.DelToken(token)
}

// 解析 token 并获得 token 过期时间, 如果返回值 err 为 ErrInvalidToken, 则需要重新验证并生成 token
func (h *Dispatcher) UnmarshalToken(token string) (user string, expiredAt int64, err error) {
	user, expiredAt, err = h.SessionDispatcher.UnmarshalToken(token)
	if session.IsInvalidTokenErr(err) {
		err = ErrInvalidToken
	}
	return
}

func (h *Dispatcher) generateAndSavePassword(id int, password string) error {
	hashedPass, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err1 != nil {
		return err1
	}
	return h.Store.SavePassword(id, string(hashedPass))
}
