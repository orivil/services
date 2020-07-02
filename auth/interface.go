// Copyright 2020 中国重庆思笃迩德人力资源有限公司. All rights reserved.
// 未经允许不可以任何形式使用

package auth

type Interface interface {
	IsOnline(token string) (ok bool, err error)
}
