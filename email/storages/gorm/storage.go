// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/orivil/services/email"
)

type Storage struct {
	db *gorm.DB
}

func (s *Storage) GetEnv() (*email.Env, error) {
	env := &email.Env{}
	err := s.db.First(env).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return nil, err
	} else {
		return env, nil
	}
}

func (s *Storage) SetEnv(env *email.Env) error {
	panic("implement me")
}
