package db

import (
	"fmt"
	"github.com/gzjjyz/srvlib/db/v2/gormx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

const MaxMysqlConnect = 8

var (
	ormEngineV2 *gorm.DB

	hasInitOrmV2 bool
	initOrmV2Mu  sync.Mutex
)

// InitOrmMysqlV2 xorm's logger and api are extremely weak, use gorm to get better feature
func InitOrmMysqlV2(user string, pwd string, host string, port int, dbs string) error {
	if hasInitOrmV2 {
		return nil
	}

	initOrmV2Mu.Lock()
	defer initOrmV2Mu.Unlock()

	if hasInitOrmV2 {
		return nil
	}

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pwd, host, port, dbs)
	if ormEngineV2, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: gormx.NewLogger(time.Second)}); err != nil {
		return err
	}

	sqlDb, err := ormEngineV2.DB()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(MaxMysqlConnect)
	sqlDb.SetMaxOpenConns(MaxMysqlConnect)

	hasInitOrmV2 = true

	return nil
}

func MustOrmMysqlV2() *gorm.DB {
	if !hasInitOrmV2 {
		panic("orm engine not init")
	}

	return ormEngineV2
}
