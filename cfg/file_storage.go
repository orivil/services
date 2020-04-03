// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import (
	"github.com/orivil/service"
	"io/ioutil"
)

type FileStorage string

type FileStorageService struct {
	storage FileStorage
	self    service.Provider
}

func NewFileStorageService(filename string) *FileStorageService {
	s := &FileStorageService{
		storage: FileStorage(filename),
		self:    nil,
	}
	s.self = s
	return s
}

func (f FileStorage) GetTomlData() ([]byte, error) {
	return ioutil.ReadFile(string(f))
}

func (f FileStorageService) New(ctn *service.Container) (value interface{}, err error) {
	return Storage(f.storage), nil
}

func (f FileStorageService) Get(ctn *service.Container) (Storage, error) {
	storage, _ := ctn.Get(&f.self)
	return storage.(Storage), nil
}
