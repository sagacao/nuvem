package dbredis

import (
	"errors"
	"io"
	"nuvem/engine/coder"
	"nuvem/engine/logger"
	"nuvem/engine/utils"
	"strings"
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

func (dbm *DBRedis) AsyncInsertRank(dbkey string, dataKey string, score uint32, uid string, rankdata coder.JSON) error {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		_, err := conn.Do("ZADD", dbkey, score, uid)
		if err != nil {
			return nil, err
		}

		savedata, err := coder.ToBytes(rankdata)
		if err != nil {
			return nil, err
		}
		_, err = conn.Do("HSET", dataKey, uid, savedata)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	_, err := exec()
	return err
}

func (dbm *DBRedis) AsyncQueryRank(dbkey string, dataKey string, start, limit uint32) ([]coder.JSON, error) {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		userscores, err := redis.StringMap(conn.Do("ZREVRANGE", dbkey, start, start+limit-1, "withscores"))
		if err != nil {
			return nil, err
		}

		replys := make([]coder.JSON, 0)
		for k, v := range userscores {
			var users []byte
			users, err = redis.Bytes(conn.Do("HGET", dataKey, k))
			if err != nil {
				if !strings.Contains(err.Error(), "nil returned") {
					logger.Error("HGET", k, err)
				}
			} else {
				jsondata := make(coder.JSON)
				err := coder.ToJSON(users, jsondata)
				if err != nil {
					logger.Error("json.Unmarshal", k, err)
				} else {
					jsondata["score"] = v
					replys = append(replys, jsondata)
				}
			}
		}
		return replys, nil
	})

	replys, err := exec()
	if replys == nil {
		return nil, err
	}
	return replys.([]coder.JSON), err
}

func (dbm *DBRedis) AsyncZAdd(dbkey string, score int64, uid string) error {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		return conn.Do("ZADD", dbkey, score, uid)
	})

	_, err := exec()
	return err
}

func (dbm *DBRedis) AsyncZRem(dbkey string, members []interface{}) error {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		args := make([]interface{}, 0)
		args = append(args, dbkey)
		for _, v := range members {
			args = append(args, v)
		}
		_, err := conn.Do("ZREM", args...)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	_, err := exec()
	return err
}

func (dbm *DBRedis) AsyncZRangeByScore(dbkey string, score1, score2 int64) ([][]byte, error) {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		return redis.ByteSlices(conn.Do("ZRANGEBYSCORE", dbkey, score1, score2))
	})

	outkeys, err := exec()
	if outkeys == nil {
		return nil, err
	}
	return outkeys.([][]byte), err
}

func (dbm *DBRedis) AsyncZRevRank(dbkey string, uid string) (int, error) {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		return redis.Int(conn.Do("ZREVRANK", dbkey, uid))
	})

	rank, err := exec()
	if err != nil {
		if strings.Contains(err.Error(), "nil") {
			return -1, nil
		}
		return -1, err
	}
	return rank.(int), err
}

func (dbm *DBRedis) AsyncHGet(dbkey string, field string) (coder.JSON, error) {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		jsonstr, err := redis.Bytes(conn.Do("HGET", dbkey, field))
		if err != nil {
			if strings.Contains(err.Error(), "nil") {
				return nil, nil
			}
			return nil, err
		}

		jsondata := make(coder.JSON)
		err = coder.ToJSON(jsonstr, jsondata)
		if err != nil {
			return nil, err
		}
		return jsondata, nil
	})

	ret, err := exec()
	if ret == nil {
		return nil, err
	}
	return ret.(coder.JSON), err
}

func (dbm *DBRedis) AsyncHGetAll(dbkey string) ([]coder.JSON, error) {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		result, err := redis.StringMap(conn.Do("HGETALL", dbkey))
		if err != nil {
			return nil, err
		}

		return result, err
	})

	ret, err := exec()
	if ret == nil {
		return nil, err
	}

	replys, ok := ret.(map[string]string)
	if !ok {
		return nil, errors.New("types error")
	}

	var res []coder.JSON
	for k, v := range replys {
		jsondata := make(coder.JSON)
		err := coder.ToJSON([]byte(v), jsondata)
		if err != nil {
			jsondata["keystr"] = k
			res = append(res, jsondata)
		}
	}
	return res, nil
}

func (dbm *DBRedis) AsyncHSet(dbkey string, field string, data coder.JSON) error {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		savedata, err := coder.ToBytes(data)
		if err != nil {
			return nil, err
		}
		_, err = conn.Do("HSET", dbkey, field, savedata)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	_, err := exec()
	return err
}

func (dbm *DBRedis) AsyncHDel(dbkey string, fields []interface{}) error {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		args := make([]interface{}, 0)
		args = append(args, dbkey)
		for _, v := range fields {
			args = append(args, v)
		}
		_, err := conn.Do("HDEL", args...)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	_, err := exec()
	return err
}

func (dbm *DBRedis) AsyncPush(cmd string, dbkey string, value []byte) error {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		return conn.Do(cmd, dbkey, value)
	})

	_, err := exec()
	return err
}

func (dbm *DBRedis) AsyncPushMulit(cmd string, dbkey string, values [][]byte) error {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		conn.Send("MULTI")
		for _, v := range values {
			conn.Send(cmd, dbkey, v)
		}
		return conn.Do("EXEC")
	})

	_, err := exec()
	return err
}

func (dbm *DBRedis) AsyncRange(cmd string, dbkey string, start, end int) ([][]byte, error) {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		return redis.ByteSlices(conn.Do(cmd, dbkey, start, end))
	})

	values, err := exec()
	if values == nil {
		return nil, err
	}
	return values.([][]byte), err
}

func (dbm *DBRedis) AsyncKeys(dbkey string) ([]string, error) {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		return redis.Strings(conn.Do("KEYS", dbkey))
	})

	ret, err := exec()
	if err != nil {
		return nil, err
	}
	return ret.([]string), err
}

func (dbm *DBRedis) AsyncDelKeys(dbkeys []interface{}) error {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		return conn.Do("DEL", dbkeys...)
	})

	_, err := exec()
	return err
}

func (dbm *DBRedis) AsyncZCard(dbkey string) (int, error) {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		return redis.Int(conn.Do("ZCARD", dbkey))
	})

	value, err := exec()
	if value == nil {
		return 0, err
	}
	return value.(int), err
}

func (dbm *DBRedis) AsyncZRank(dbkey string, uid string) (int, error) {
	exec := utils.Future(func() (interface{}, error) {
		conn := dbm.pool.Get()
		defer conn.Close()

		return redis.Int(conn.Do("ZRANK", dbkey, uid))
	})

	value, err := exec()
	if value == nil {
		return 0, err
	}
	return value.(int), err
}
