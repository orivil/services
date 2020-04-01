// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import (
	"github.com/orivil/service"
	"github.com/orivil/xcfg"
)

var defFile = "configs/config.toml"

type Service struct {
	file string
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var env xcfg.Env
	env, err = xcfg.DecodeFile(s.file)
	if err != nil {
		return nil, err
	}
	return env, nil
}

func NewService(configFile string) *Service {
	if configFile == "" {
		configFile = defFile
	}
	return &Service{file: configFile}
}
