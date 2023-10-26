package database

import (
	"fmt"
	"strconv"

	"github.com/bookpanda/mygraderlist-auth/src/app/model/auth"
	"github.com/bookpanda/mygraderlist-auth/src/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(conf *config.Database) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", conf.User, conf.Password, conf.Host, strconv.Itoa(conf.Port), conf.Name)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(auth.Auth{})
	if err != nil {
		return nil, err
	}

	return
}
