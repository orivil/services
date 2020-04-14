// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package email

import (
	"bytes"
	"mime/quotedprintable"
	"net/smtp"
)

type Sender struct {
	FromEmail string
	Auth      smtp.Auth
	Addr      string
}

func NewSMTPSender(env *Env) *Sender {
	return &Sender{
		FromEmail: env.From,
		Auth:      smtp.PlainAuth("", env.Username, env.Password, env.Host),
		Addr:      env.Host + env.Port,
	}
}

func (et *Sender) Send(toEmails []string, subject, contentType string, body []byte) error {
	body = InitEmailBody(et.FromEmail, subject, contentType, body)
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
