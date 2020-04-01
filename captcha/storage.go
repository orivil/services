// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha

import "time"

type Storage interface {
	GetCaptcha(id string) (string, error)
	SetCaptcha(id, captcha string, expires time.Duration) error
	IsCaptchaOK(id, captcha string) (ok bool, err error)
	DelCaptcha(id string) (err error)
}
