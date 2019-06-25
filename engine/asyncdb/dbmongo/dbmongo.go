package dbmongo

import (
	"container/heap"
	"io"
	"nuvem/engine/coder"
	"nuvem/engine/logger"
	"nuvem/engine/utils"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoSession struct {
	*mgo.Session
	ref   int
	index int
}

type SessionHeap []*MongoSession

func (h SessionHeap) Len() int {
	return len(h)
}

func (h SessionHeap) Less(i, j int) bool {
	return h[i].ref < h[j].ref
}

func (h SessionHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *SessionHeap) Push(s interface{}) {
	s.(*MongoSession).index = len(*h)
	*h = append(*h, s.(*MongoSession))
}

func (h *SessionHeap) Pop() interface{} {
	l := len(*h)
	s := (*h)[l-1]
	s.index = -1
	*h = (*h)[:l-1]
	return s
}

type DBMongo struct {
	sync.Mutex
	driverName string
	dataSource string
	heaps      SessionHeap
}

// OpenMongo opens SQL driver for KVDB backend
func OpenMongo(url string, dbname string, collectionName string) (*DBMongo, error) {
	db, err := OpenMongoWithTimeout(url, 10, 10*time.Second, 5*time.Minute)
	if err == nil {
		logger.Info("Mongo Has Opened at ", url)
	}
	db.dataSource = dbname
	logger.Debug("Mongo:Dial ", url, ":", dbname, " success!!!")
	return db, err
}

// OpenMongoWithTimeout opens SQL driver for KVDB backend .goroutine safe
func OpenMongoWithTimeout(url string, sessionNum int, dialTimeout time.Duration, timeout time.Duration) (*DBMongo, error) {
	if sessionNum <= 0 {
		sessionNum = 100
		logger.Debug("invalid sessionNum, reset to ", sessionNum)
	}

	s, err := mgo.DialWithTimeout(url, dialTimeout)
	if err != nil {
		return nil, err
	}
	s.SetSyncTimeout(timeout)
	s.SetSocketTimeout(timeout)

	// SessionHeap
	heaps := make(SessionHeap, sessionNum)
	heaps[0] = &MongoSession{s, 0, 0}
	for i := 1; i < sessionNum; i++ {
		heaps[i] = &MongoSession{s.New(), 0, i}
	}
	heap.Init(&heaps)

	return &DBMongo{
		driverName: "mongo",
		heaps:      heaps,
	}, nil
}

//Close goroutine safe
func (dbm *DBMongo) Close() {
	dbm.Lock()
	for _, s := range dbm.heaps {
		s.Close()
		if s.ref != 0 {
			logger.Error("session ref = ", s.ref)
		}
	}
	dbm.Unlock()
}

//Ref goroutine safe
func (dbm *DBMongo) Ref() *MongoSession {
	dbm.Lock()
	s := dbm.heaps[0]
	if s.ref == 0 {
		s.Refresh()
	}
	s.ref++
	heap.Fix(&dbm.heaps, 0)
	dbm.Unlock()

	return s
}

//UnRef goroutine safe
func (dbm *DBMongo) UnRef(s *MongoSession) {
	dbm.Lock()
	s.ref--
	heap.Fix(&dbm.heaps, s.index)
	dbm.Unlock()
}

func (dbm *DBMongo) IsConnectionError(err error) bool {
	return err == io.EOF || err == io.ErrUnexpectedEOF
}

func (dbm *DBMongo) AsyncQuery(collection string, query interface{}) (*mgo.Iter, error) {
	exec := utils.Future(func() (interface{}, error) {
		s := dbm.Ref()
		defer dbm.UnRef(s)

		reply := s.DB(dbm.dataSource).C(collection).Find(query).Iter()
		return reply, reply.Err()
	})

	ret, err := exec()
	if ret == nil {
		return nil, err
	}
	return ret.(*mgo.Iter), err
}

func (dbm *DBMongo) AsyncExec(collection string, dbkey string, data interface{}) (*mgo.ChangeInfo, error) {
	exec := utils.Future(func() (interface{}, error) {
		s := dbm.Ref()
		defer dbm.UnRef(s)

		if dbkey == "" {
			err := s.DB(dbm.dataSource).C(collection).Insert(data)
			return nil, err
		}
		return s.DB(dbm.dataSource).C(collection).UpsertId(dbkey, data)
	})

	ret, err := exec()
	if ret == nil {
		return nil, err
	}
	return ret.(*mgo.ChangeInfo), err
}

func (dbm *DBMongo) AsyncUpdate(collection string, query bson.M, data bson.M) (*mgo.ChangeInfo, error) {
	exec := utils.Future(func() (interface{}, error) {
		s := dbm.Ref()
		defer dbm.UnRef(s)

		return s.DB(dbm.dataSource).C(collection).Upsert(query, bson.M{"$set": data})
	})

	ret, err := exec()
	if ret == nil {
		return nil, err
	}
	return ret.(*mgo.ChangeInfo), err
}

func (dbm *DBMongo) AsyncQueryDB(dbname, collection string, query interface{}) (*mgo.Iter, error) {
	exec := utils.Future(func() (interface{}, error) {
		s := dbm.Ref()
		defer dbm.UnRef(s)

		reply := s.DB(dbname).C(collection).Find(query).Iter()
		return reply, reply.Err()
	})

	ret, err := exec()
	if ret == nil {
		return nil, err
	}
	return ret.(*mgo.Iter), err
}

func (dbm *DBMongo) AsyncUpdateDB(dbname, collection string, query bson.M, data coder.JSON) error {
	exec := utils.Future(func() (interface{}, error) {
		s := dbm.Ref()
		defer dbm.UnRef(s)

		err := s.DB(dbname).C(collection).Update(query, bson.M{"$set": data})
		return nil, err
	})

	_, err := exec()
	return err
}
