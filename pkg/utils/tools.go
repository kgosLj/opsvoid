package utils

import (
	"fmt"
	"github.com/kgosLj/opsvoid/config"
)

func GetDSN(db config.MysqlConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.DBName)
}
