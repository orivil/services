// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package gorm

import (
	"github.com/jinzhu/gorm"
	"github.com/orivil/service"
	"github.com/orivil/services/cfg"
	"github.com/orivil/xcfg"
)

var Service service.Provider = func(c *service.Container) (value interface{}, err error) {
	env := &Env{}
	configs := c.MustGet(&cfg.Provider).(xcfg.Data)
	err = configs.Unmarshal(env)
	if err != nil {
		return nil, err
	}
	var gormDB *gorm.DB
	gormDB, err = env.Connect()
	if err != nil {
		return nil, err
	} else {
		return gormDB, nil
	}
}
