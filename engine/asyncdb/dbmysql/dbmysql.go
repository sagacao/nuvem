package dbmysql

import (
	"database/sql"
	"fmt"
	"nuvem/engine/logger"
	"nuvem/engine/utils"

	_ "github.com/go-sql-driver/mysql"
)

const (
	_MAX_KEY_LENGTH = 256
)

type DBMysql struct {
	driverName string
	dataSource string
	db         *sql.DB
}

// OpenMySQL opens SQL driver for KVDB backend
func OpenMySQL(dataSource string) (*DBMysql, error) {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)

	return &DBMysql{
		driverName: "mysql",
		dataSource: dataSource,
		db:         db,
	}, nil
}

func (dbm *DBMysql) String() string {
	return fmt.Sprintf("%s<%s>", dbm.driverName, dbm.dataSource)
}

func (dbm *DBMysql) AsyncQuery(query string, args ...interface{}) (*sql.Rows, error) {
	exec := utils.Future(func() (interface{}, error) {
		return dbm.db.Query(query, args...)
	})

	ret, err := exec()
	if ret == nil {
		return nil, err
	}
	return ret.(*sql.Rows), err
}

func (dbm *DBMysql) AsyncExec(query string, args ...interface{}) (sql.Result, error) {
	exec := utils.Future(func() (interface{}, error) {
		return dbm.db.Exec(query, args...)
	})

	ret, err := exec()
	if err == nil {
		return nil, err
	}
	return ret.(sql.Result), err
}

func (dbm *DBMysql) Close() {
	if err := dbm.db.Close(); err != nil {
		logger.Error("%s: close error: %s", dbm.String(), err)
	}
}

func (dbm *DBMysql) IsConnectionError(err error) bool {
	return true
}
