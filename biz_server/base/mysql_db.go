package base

import (
	"database/sql"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

var MysqlDB *sql.DB

func init() {
	var mysqlErr error
	MysqlDB, mysqlErr = sql.Open("mysql", "root:123456@@tcp(111.230.244.177:3306)/hero_story")

	if nil != mysqlErr {
		panic(mysqlErr)
	}

	MysqlDB.SetMaxOpenConns(128)
	MysqlDB.SetMaxIdleConns(16)
	MysqlDB.SetConnMaxLifetime(2 * time.Minute)

	if mysqlErr = MysqlDB.Ping(); nil != mysqlErr {
		panic(mysqlErr)
	}
}
