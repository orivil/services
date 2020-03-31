// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Env struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

func (e Env) Init(db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     e.Addr,
		Password: e.Password,
		DB:       db,
	})
	err := client.Ping().Err()
	if err != nil {
		return nil, fmt.Errorf("ping redis: %s", err)
	} else {
		return client, nil
	}
}
