package dbredis

import (
	"io"
	"nuvem/engine/logger"
	"time"

	"github.com/garyburd/redigo/redis"
)

type DBRedis struct {
	driverName string
	dataSource string
	pool       *redis.Pool
}

// OpenRedis opens SQL driver for KVDB backend
func OpenRedis(url string, prefix string, passwd string, dbindex int) (*DBRedis, error) {
	db := &DBRedis{
		driverName: "redis",
		dataSource: url,
	}

	db.pool = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   30,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				return nil, err
			}
			if passwd != "" {
				if _, err := c.Do("AUTH", passwd); err != nil {
					c.Close()
					return nil, err
				}
			}

			if dbindex > 0 {
				if _, err := c.Do("SELECT", dbindex); err != nil {
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	logger.Info("Redis Has Open at ", url)
	return db, nil
}

func (dbm *DBRedis) Close() {
	if dbm.pool != nil {
		dbm.pool.Close()
	}
}

func (dbm *DBRedis) IsConnectionError(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}

func (dbm *DBRedis) AsyncHGet() {

}
