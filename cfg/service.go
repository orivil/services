// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import (
	"github.com/orivil/service"
	"github.com/orivil/xcfg"
)

var ConfigFile = "configs/config.toml"

var Service service.Provider = func(c *service.Container) (value interface{}, err error) {
	var env xcfg.Env
	env, err = xcfg.DecodeFile(ConfigFile)
	if err != nil {
		return nil, err
	}
	return env, nil
}
