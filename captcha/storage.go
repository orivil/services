// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha

import (
	"github.com/orivil/service"
	"time"
)

type Storage interface {
	SetCaptcha(id, captcha string, expires time.Duration) error
	IsCaptchaOK(id, captcha string) (ok bool, err error)
	DelCaptcha(id string) (err error)
}

type StorageService interface {
	service.Provider
	Get(ctn *service.Container) (Storage, error)
}
