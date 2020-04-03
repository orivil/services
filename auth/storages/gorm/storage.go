// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package auth_gorm_storage

import "github.com/jinzhu/gorm"

type Storage struct {
	db *gorm.DB
}

func (s *Storage) GetPassword(id int) (password string, err error) {
	up := &UserPassword{}
	err = s.db.Where("us_id=?", id).First(up).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return up.Pass, nil
}

// TODO: 测试主键不是ID是否被保存
func (s *Storage) SavePassword(id int, password string) error {
	up := &UserPassword{
		UsID: id,
		Pass: password,
	}
	return s.db.Save(up).Error
}

func NewStorage(db *gorm.DB, migrate bool) *Storage {
	if migrate {
		db.AutoMigrate(&UserPassword{})
	}
	return &Storage{
		db: db,
	}
}
