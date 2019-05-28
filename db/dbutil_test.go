package db

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

func TestDbTrans(t *testing.T) {
	dsn := "med_dev:&*nXbL^LDxc0@tcp(118.31.236.23:3306)/?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	db, err := gorm.Open("mysql", dsn)
	defer db.Close()

	db_func := func(db *gorm.DB) error {
		test_sql := "update med_medication.medication set otc_type = 11 where id = 295"
		eerr := db.Exec(test_sql).Error
		if eerr != nil {
			return errors.Wrap(eerr, "更新295失败")
		}

		test_sql = "update med_medication.medication set otc_type = 22 where id = 296"
		eerr = db.Exec(test_sql).Error
		if eerr != nil {
			return errors.Wrap(eerr, "更新296失败")
		}

		return nil
	}

	err = DoTrans(db, db_func)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("test ok")
}

func TestDbTransPanic(t *testing.T) {
	dsn := "med_dev:&*nXbL^LDxc0@tcp(118.31.236.23:3306)/?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	db, err := gorm.Open("mysql", dsn)
	defer db.Close()

	db_func := func(db *gorm.DB) error {
		test_sql := "update med_medication.medication set otc_type = 11 where id = 295"
		eerr := db.Exec(test_sql).Error
		if eerr != nil {
			return errors.Wrap(eerr, "更新295失败")
		}

		val := 0
		t.Log(100 / val)

		return nil
	}

	err = DoTrans(db, db_func)
	if err == nil {
		t.Fatal("未返回错误")
	}
	t.Log(err.Error())
}
