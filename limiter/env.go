// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package limiter

import (
	"github.com/orivil/limiter"
	"time"
)

/**
# 失败次数限制器
[failed-times-limiter]

# 等待时间(单位: 秒)
wait = 60

# 第几次开始等待, 达到等待次数后等待时间呈指数增加, 直到等待过期
# 或失败次数被清零
start_limit_times = 5
*/

type Env struct {
	Wait            int64 `toml:"wait"`
	StartLimitTimes int64 `toml:"start_limit_times"`
}

func toLimiterOptions(e *Env) *limiter.Options {
	return &limiter.Options{
		Wait:            time.Duration(e.Wait) * time.Second,
		StartLimitTimes: e.StartLimitTimes,
	}
}
