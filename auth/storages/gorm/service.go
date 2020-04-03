// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package auth_gorm_storage

import (
	"github.com/orivil/service"
	"github.com/orivil/services/auth"
	"github.com/orivil/services/database/gorm"
)

type Service struct {
	gormService *gorm.Service
	self        service.Provider
}

func (s *Service) New(ctn *service.Container) (interface{}, error) {
	db, err := s.gormService.Get(ctn)
	if err != nil {
		return nil, err
	}
	return auth.Storage(NewStorage(db, true)), nil
}

func (s *Service) Get(ctn *service.Container) (auth.Storage, error) {
	store, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return store.(auth.Storage), nil
	}
}

func NewService(gormService *gorm.Service) *Service {
	s := &Service{
		gormService: gormService,
		self:        nil,
	}
	s.self = s
	return s
}
