// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session_memory_storage

import (
	"sync"
	"time"
)

var NowFunc = time.Now

type Storage struct {
	sessions map[string]*time.Time
	mu       sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{sessions: map[string]*time.Time{}}
}

func (s *Storage) IsOK(session string) (ok bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var exp *time.Time
	exp, ok = s.sessions[session]
	if ok && exp.After(NowFunc()) {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *Storage) SaveSession(session string, expires time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp := NowFunc().Add(expires)
	s.sessions[session] = &exp
	return nil
}

func (s *Storage) DelSession(session string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, session)
	return nil
}
