// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email

import (
	"bytes"
	"html/template"
	"mime/quotedprintable"
	"net/smtp"
	"sync"
)

type Sender struct {
	FromEmail string
	Auth      smtp.Auth
	tpl       *template.Template
	Addr      string
}

func NewSMTPSender(env *Env, tplStorage TemplateStorage) (*Sender, error) {
	data, err := tplStorage.Read()
	if err != nil {
		return nil, err
	}
	var tpl *template.Template
	tpl, err = template.New("email").Parse(string(data))
	if err != nil {
		return nil, err
	}
	return &Sender{
		FromEmail: env.From,
		tpl:       tpl,
		Auth:      smtp.PlainAuth("", env.Username, env.Password, env.Host),
		Addr:      env.Host + env.Port,
	}, nil
}

var bodyPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func (et *Sender) Send(toEmails []string, subject, contentType string, data interface{}) error {
	buf := bodyPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bodyPool.Put(buf)
	err := et.tpl.Execute(buf, data)
	if err != nil {
		return err
	}
	body := InitEmailBody(et.FromEmail, subject, contentType, buf.Bytes())
	return smtp.SendMail(et.Addr, et.Auth, et.FromEmail, toEmails, body)
}

func InitEmailBody(from, subject, contentType string, body []byte) []byte {
	msg := []byte("Subject: " + subject + "\r\n" +
		"From: " + from + "\r\n" +
		"MIME-Version: 1.0" + "\r\n" +
		"Content-Type: " + contentType + "\r\n" + // text/html; charset=UTF-8
		"Content-Transfer-Encoding: quoted-printable" + "\r\n" +
		"Content-Disposition: inline" + "\r\n")
	buffer := bytes.NewBuffer(msg)
	w := quotedprintable.NewWriter(buffer)
	_, _ = w.Write(body)
	return buffer.Bytes()
}
