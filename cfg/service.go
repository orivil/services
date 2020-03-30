// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cfg

import (
	"github.com/orivil/service"
	"github.com/orivil/xcfg"
)

var Provider service.Provider = func(c *service.Container) (value interface{}, err error) {
	return make(xcfg.Data), nil
}

func Init(configs xcfg.Data, container *service.Container) {
	container.SetCache(&Provider, configs)
}
