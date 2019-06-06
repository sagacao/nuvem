package bridge

import (
	"errors"
	"nuvem/engine/coder"
	"nuvem/engine/logger"
	"sync"
)

// eventConn wraps Conn with a specified event type.
type clientConn struct {
	uid  string
	sock *Socket
}

// binder is defined to store the relation of userID and eventConn
type binder struct {
	mu sync.RWMutex

	// map stores key: connID and value: *clientConn
	connID2Sockets map[string]*clientConn
}

// Bind binds userID with eConn specified by event. It fails if the
// return error is not nil.
func (b *binder) Bind(userID string, conn *Socket) error {
	if userID == "" {
		return errors.New("userID can't be empty")
	}

	if conn == nil {
		return errors.New("conn can't be nil")
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.connID2Sockets[conn.GetID()] = &clientConn{uid: userID, sock: conn}

	return nil
}

// Unbind unbind and removes Conn if it's exist.
func (b *binder) Unbind(conn *Socket) error {
	if conn == nil {
		return errors.New("conn can't be empty")
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	//TODO
	msg, err := coder.ToBytes(coder.JSON{"mid": 1002})
	if err != nil {
		logger.Error("OnMessage", conn.GetID(), err)
	} else {
		conn.Write(msg)
	}

	delete(b.connID2Sockets, conn.GetID())
	return nil
}

func (b *binder) ForEach(f func(sock *Socket)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, s := range b.connID2Sockets {
		f(s.sock)
	}
}

// FindConn trys to find Conn by ID.
func (b *binder) FindConn(connID string) (*Socket, bool) {
	if connID == "" {
		return nil, false
	}

	client, ok := b.connID2Sockets[connID]
	// if userID been found by connID, then find the Conn using userID
	if ok {
		return client.sock, true
	}

	// // userID not found, iterate all the conns
	// for _, eConns := range b.userID2EventConnMap {
	// 	for i := range *eConns {
	// 		if (*eConns)[i].Sock.GetID() == connID {
	// 			return (*eConns)[i].Sock, true
	// 		}
	// 	}
	// }

	return nil, false
}

// FilterConn searches the conns related to userID, and filtered by
// event. The userID can't be empty. The event will be ignored if it's empty.
// All the conns related to the userID will be returned if the event is empty.
func (b *binder) FilterConn(userID, event string) ([]*Socket, error) {
	if userID == "" {
		return nil, errors.New("userID can't be empty")
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	// if eConns, ok := b.userID2EventConnMap[userID]; ok {
	// 	ecs := make([]*Socket, 0, len(*eConns))
	// 	for i := range *eConns {
	// 		if event == "" || (*eConns)[i].Event == event {
	// 			ecs = append(ecs, (*eConns)[i].Sock)
	// 		}
	// 	}
	// 	return ecs, nil
	// }

	return []*Socket{}, nil
}

func (b *binder) BindCount() int {
	b.mu.RLock()
	bindCount := len(b.connID2Sockets)
	b.mu.RUnlock()
	return bindCount
}
