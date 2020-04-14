// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email_captcha

import (
	"bytes"
	"github.com/orivil/services/captcha"
	"github.com/orivil/services/email"
	"strings"
	"sync"
	"text/template"
	"time"
)

var ExecuteBody = func(tpl *template.Template, buf *bytes.Buffer, code string) error {
	return tpl.Execute(buf, code)
}

type Dispatcher struct {
	emailSender *email.Sender
	store       captcha.Storage
	tpl         *template.Template
	env         *Env
}

func NewDispatcher(env *Env, sender *email.Sender, storage captcha.Storage, tpl *template.Template) *Dispatcher {
	return &Dispatcher{
		emailSender: sender,
		store:       storage,
		tpl:         tpl,
		env:         env,
	}
}

var bodyPool = sync.Pool{New: func() interface{} {
	return new(bytes.Buffer)
}}

func (d *Dispatcher) SendCaptcha(subject, email string) error {
	code := d.env.GenCaptcha()
	err := d.store.SetCaptcha(email, code, time.Duration(d.env.Expires)*time.Second)
	if err != nil {
		return err
	}
	buf := bodyPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		bodyPool.Put(buf)
	}()
	err = ExecuteBody(d.tpl, buf, code)
	if err != nil {
		return err
	}
	return d.emailSender.Send([]string{email}, subject, d.env.ContentType, buf.Bytes())
}

func (d *Dispatcher) Verify(email, code string) (ok bool, err error) {
	return d.store.IsCaptchaOK(email, strings.ToUpper(code))
}

func (d *Dispatcher) Remove(email string) error {
	return d.store.DelCaptcha(email)
}
