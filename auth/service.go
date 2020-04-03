// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package auth

import (
	"github.com/orivil/service"
	"github.com/orivil/services/session"
)

type Service struct {
	storageService StorageService
	sessionService *session.Service
	self           service.Provider
}

func NewService(sessionService *session.Service, storageService StorageService) *Service {
	s := &Service{
		storageService: storageService,
		sessionService: sessionService,
		self:           nil,
	}
	s.self = s
	return s
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var store Storage
	store, err = s.storageService.Get(ctn)
	if err != nil {
		return nil, err
	}
	var dispatcher *session.Dispatcher
	dispatcher, err = s.sessionService.Get(ctn)
	if err != nil {
		return nil, err
	}
	return NewDispatcher(store, dispatcher), nil
}

func (s *Service) Get(ctn *service.Container) (*Dispatcher, error) {
	dis, err := ctn.Get(&s.self)
	if err != nil {
		return nil, err
	} else {
		return dis.(*Dispatcher), nil
	}
}
