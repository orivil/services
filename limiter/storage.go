// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package limiter

import (
	"github.com/orivil/limiter"
	"github.com/orivil/service"
)

type Storage interface {
	limiter.Storage
}

type StorageService interface {
	service.Provider
	Get(container *service.Container) (Storage, error)
}

type MemoryStorageService struct {
	self service.Provider
}

func NewMemoryStorageService() *MemoryStorageService {
	s := &MemoryStorageService{
		self: nil,
	}
	s.self = s
	return s
}

func (m *MemoryStorageService) New(ctn *service.Container) (value interface{}, err error) {
	return Storage(limiter.NewMemoryStorage()), nil
}

func (m *MemoryStorageService) Get(ctn *service.Container) (Storage, error) {
	v, err := ctn.Get(&m.self)
	if err != nil {
		return nil, err
	}
	return v.(Storage), nil
}
