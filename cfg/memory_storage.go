// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import (
	"github.com/orivil/service"
)

type MemoryStorage string

type MemoryStorageService struct {
	storage MemoryStorage
	self    service.Provider
}

func NewMemoryStorageService(configData string) *MemoryStorageService {
	s := &MemoryStorageService{
		storage: MemoryStorage(configData),
		self:    nil,
	}
	s.self = s
	return s
}

func (s MemoryStorage) GetTomlData() ([]byte, error) {
	return []byte(s), nil
}

func (m MemoryStorageService) New(ctn *service.Container) (value interface{}, err error) {
	return Storage(m.storage), nil
}

func (m MemoryStorageService) Get(ctn *service.Container) (Storage, error) {
	storage, _ := ctn.Get(&m.self)
	return storage.(Storage), nil
}
