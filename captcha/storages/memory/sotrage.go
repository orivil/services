// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package captcha_memory_storage

import (
	"sync"
	"time"
)

var NowFunc = time.Now

type Storage struct {
	vs map[string]*exValue
	mu sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		vs: make(map[string]*exValue, 10),
		mu: sync.Mutex{},
	}
}

type exValue struct {
	expireAt *time.Time
	value    string
}

func (s *Storage) getCaptcha(id string) (string, error) {
	if v, ok := s.vs[id]; ok {
		return v.value, nil
	} else {
		return "", nil
	}
}

func (s *Storage) SetCaptcha(id, captcha string, expires time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := NowFunc()
	// 随机取 3 个数据查看是否过期, 过期则删除
	max := 3
	for k, value := range s.vs {
		if value.expireAt.Before(now) {
			delete(s.vs, k)
		}
		max--
		if max <= 1 {
			break
		}
	}
	exAt := now.Add(expires)
	v := &exValue{
		expireAt: &exAt,
		value:    captcha,
	}
	s.vs[id] = v
	return nil
}

func (s *Storage) IsCaptchaOK(id, captcha string) (ok bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var exist string
	exist, err = s.getCaptcha(id)
	if err != nil {
		return false, err
	} else {
		return exist == captcha, nil
	}
}

func (s *Storage) DelCaptcha(id string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.vs, id)
	return nil
}
