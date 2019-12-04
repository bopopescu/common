package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Task func(db *gorm.DB) error

func DoTrans(db *gorm.DB, tasks ...Task) (rerr error) {
	new_db := db.Begin()
	if err := new_db.Error; err != nil {
		return errors.Wrap(err, "启动事务失败")
	}

	defer func() {
		if e := recover(); e != nil {
			new_db.Rollback()
			rerr = errors.Errorf("单个事务发生崩溃 %s", e.Error())
		}
	}()

	for _, task := range tasks {
		if err := task(new_db); err != nil {
			if rerr := new_db.Rollback().Error; rerr != nil {
				return errors.Wrap(rerr, "回滚事务失败")
			}
			return err
		}
	}

	if err := new_db.Commit().Error; err != nil {
		new_db.Rollback()
		return errors.Wrap(err, "提交事务失败")
	}
	return nil
}
