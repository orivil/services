// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package image_captcha

import (
	"github.com/orivil/captcha"
	captcha2 "github.com/orivil/services/captcha"
	"io"
	"strings"
	"time"
)

// ID 长度
const IDLength = 12

type Dispatcher struct {
	store captcha2.Storage
	env   *Env
}

func NewDispatcher(store captcha2.Storage, env *Env) *Dispatcher {
	return &Dispatcher{
		store: store,
		env:   env,
	}
}

func (s *Dispatcher) GenID() string {
	return string(captcha2.NewID(IDLength))
}

// 生成验证码并将图片写入接口中
func (s *Dispatcher) WriteCodeToImage(id string, w io.Writer) error {
	code := s.env.GenCaptcha()
	err := s.store.SetCaptcha(id, code, time.Duration(s.env.Expires)*time.Second)
	if err != nil {
		return err
	}
	_, err = captcha.NewImage(id, []byte(code), s.env.ImgWidth, s.env.ImgHeight).WriteTo(w)
	return err
}

// 验证验证码，如果验证失败，则需要重新生成验证码
func (s *Dispatcher) Verify(id, code string) (ok bool, err error) {
	ok, err = s.store.IsCaptchaOK(id, strings.ToUpper(code))
	if err != nil {
		return false, err
	} else {
		err = s.store.DelCaptcha(id)
		if err != nil {
			return false, err
		} else {
			return ok, nil
		}
	}
}
