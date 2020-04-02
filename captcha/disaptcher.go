// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha

import (
	"github.com/orivil/captcha"
	"io"
	"math/rand"
	"strings"
	"time"
)

// ID 长度
const IDLength = 16

// ID 字符集
var idChars = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 验证码字符集
var numberChars = []byte("123567890")                                 // 无 4, 0
var letterChars = []byte("ABCDEFGHJKLMNPQRSTUVWXYZ")                  // 无 i, O
var numberAndLetterChars = []byte("12356789ABCDEFGHJKLMNPQRSTUVWXYZ") // 无 4, 0, i, O

// 从 src 中随机创建 bit 位的随机字符数组
func randomBytes(src []byte, bit int) []byte {
	bts := make([]byte, bit)
	r := make([]byte, bit)
	ln := len(src)
	rand.Read(bts)
	for idx, b := range bts {
		r[idx] = src[int(b)%ln]
	}
	return r
}

type Dispatcher struct {
	store Storage
	env   *Env
}

func NewDispatcher(store Storage, env *Env) *Dispatcher {
	return &Dispatcher{
		store: store,
		env:   env,
	}
}

func (s *Dispatcher) GenID() (id string, err error) {
	id = string(randomBytes(idChars, IDLength))
	var chars []byte
	switch s.env.Type {
	case Number:
		chars = numberChars
	case Letter:
		chars = letterChars
	case LetterAndNumber:
		chars = numberAndLetterChars
	}
	code := strings.ToUpper(string(randomBytes(chars, s.env.CodeLength)))
	err = s.store.SetCaptcha(id, code, time.Duration(s.env.Expires)*time.Second)
	if err != nil {
		return "", err
	} else {
		return id, nil
	}
}

func (s *Dispatcher) WriteImage(id string, w io.Writer) error {
	code, err := s.store.GetCaptcha(id)
	if err != nil {
		return err
	}
	_, err = captcha.NewImage(id, []byte(code), s.env.ImgWidth, s.env.ImgHeight).WriteTo(w)
	return err
}

// 验证验证码，如果验证失败，则需要重新生成验证码
func (s *Dispatcher) Verify(id, code string) (ok bool, err error) {
	ok, err = s.store.IsCaptchaOK(id, code)
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
