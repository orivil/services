// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"fmt"
	"github.com/orivil/service"
)

type Dispatcher struct {
}

type Service struct {
	self service.Provider
}

func NewService() *Service {
	s := &Service{}
	s.self = s
	return s
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	return &Dispatcher{}, nil
}

func (s *Service) Get(ctn *service.Container) (dispatcher *Dispatcher, err error) {
	d, er := ctn.Get(&s.self)
	if er != nil {
		return nil, er
	} else {
		return d.(*Dispatcher), nil
	}
}

func main() {
	var s interface{} = &Service{}
	_, ok := s.(service.Provider)
	fmt.Println(ok)
}
