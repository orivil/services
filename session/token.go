// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package session

type Token struct {
	ID           string
	RefreshToken string
	ExpiresAt    int64
}
