package bridge

import (
	"errors"
	"sync"
	"time"

	"nuvem/engine/logger"
	"nuvem/engine/uuid"

	"github.com/gorilla/websocket"
)

// Conn wraps websocket.Conn with Conn. It defines to listen and read
// data from Conn.
type Socket struct {
	Conn *websocket.Conn

	AfterReadFunc   func(messageType int, rmessage []byte)
	BeforeCloseFunc func()

	once   sync.Once
	id     string
	stopCh chan struct{}
}

// Write write p to the websocket connection. The error returned will always
// be nil if success.
func (c *Socket) Write(p []byte) (n int, err error) {
	select {
	case <-c.stopCh:
		return 0, errors.New("Socket is closed, can't be written")
	default:
		err = c.Conn.WriteMessage(websocket.BinaryMessage, p)
		if err != nil {
			return 0, err
		}
		return len(p), nil
	}
}

// GetID returns the id generated using UUID algorithm.
func (c *Socket) GetID() string {
	c.once.Do(func() {
		c.id = uuid.GenUUID()
	})

	return c.id
}

// Listen listens for receive data from websocket connection. It blocks
// until websocket connection is closed.
func (c *Socket) Listen() {
	c.Conn.SetCloseHandler(func(code int, text string) error {
		if c.BeforeCloseFunc != nil {
			c.BeforeCloseFunc()
		}

		if err := c.Close(); err != nil {
			logger.Error(err)
		}

		message := websocket.FormatCloseMessage(code, "")
		c.Conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		return nil
	})

	// Keeps reading from Conn util get error.
ReadLoop:
	for {
		select {
		case <-c.stopCh:
			break ReadLoop
		default:
			messageType, r, err := c.Conn.ReadMessage() //c.Conn.NextReader()
			if err != nil {
				// TODO: handle read error maybe
				break ReadLoop
			}

			if c.AfterReadFunc != nil {
				c.AfterReadFunc(messageType, r)
			}
		}
	}
}

// Close close the connection.
func (c *Socket) Close() error {
	select {
	case <-c.stopCh:
		return errors.New("Socket already been closed")
	default:
		c.Conn.Close()
		close(c.stopCh)
		return nil
	}
}

// NewSocket wraps conn.
func NewSocket(conn *websocket.Conn) *Socket {
	return &Socket{
		Conn:   conn,
		stopCh: make(chan struct{}),
	}
}
