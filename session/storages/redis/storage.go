// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session_redis_storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{client: client}
}

func (s *Storage) IsOK(session string) (ok bool, err error) {
	res, er := s.client.Exists(context.Background(), session).Result()
	if er != nil {
		return false, er
	}
	return res == 1, nil
}

func (s *Storage) SaveSession(session string, expires time.Duration) error {
	return s.client.Set(context.Background(), session, 1, expires).Err()
}

func (s *Storage) DelSession(session string) error {
	return s.client.Del(context.Background(), session).Err()
}
