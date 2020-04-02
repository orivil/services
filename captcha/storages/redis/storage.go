// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package redis

import (
	"github.com/go-redis/redis"
	"time"
)

type Storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{client: client}
}

func (c *Storage) GetCaptcha(id string) (string, error) {
	captcha, err := c.client.Get(id).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
	}
	return captcha, err
}

func (c *Storage) SetCaptcha(id, captcha string, expires time.Duration) error {
	return c.client.Set(id, captcha, expires).Err()
}

func (c *Storage) IsCaptchaOK(id, captcha string) (ok bool, err error) {
	exist, err1 := c.client.Get(id).Result()
	if err1 != nil {
		if err1 == redis.Nil {
			return false, nil
		} else {
			return false, err1
		}
	} else {
		return exist == captcha, nil
	}
}

func (c *Storage) DelCaptcha(id string) (err error) {
	return c.client.Del(id).Err()
}