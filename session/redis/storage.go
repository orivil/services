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

func (s *Storage) GetOnlineSessions(id string) (total int, err error) {
	i64, err1 := s.client.SCard(s.idKey(id)).Result()
	if err1 != nil {
		return 0, err1
	}
	return int(i64), nil
}

func (s *Storage) idKey(id string) string {
	return "i_k_" + id
}

func (s *Storage) sKey(session string) string {
	return "s_k_" + session
}

func (s *Storage) IsOK(id, session string) (ok bool, err error) {
	i, err1 := s.client.Exists(s.sKey(session)).Result()
	if err1 != nil {
		return false, err1
	}
	return i == 1, nil
}

func (s *Storage) SaveSession(id, session string, expires time.Duration) error {
	return s.client.SAdd(id, session).Err()
}

func (s *Storage) DelSession(id, session string) error {
	return s.client
}

func (s *Storage) DelFirstExpireSession(id string) error {
	panic("implement me")
}
