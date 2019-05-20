package asura

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Session wrapper around websocket connections.
type Socket struct {
	id      string
	Keys    map[string]interface{}
	conn    *websocket.Conn
	output  chan *envelope
	asura   *Asura
	open    bool
	rwmutex *sync.RWMutex
}

func (s *Socket) Sid() string {
	return s.id
}

func (s *Socket) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *Socket) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Socket) writeMessage(message *envelope) {
	if s.closed() {
		s.asura.errorHandler(s, errors.New("tried to write to closed a session"))
		return
	}

	select {
	case s.output <- message:
	default:
		s.asura.errorHandler(s, errors.New("session message buffer is full"))
	}
}

func (s *Socket) writeRaw(message *envelope) error {
	if s.closed() {
		return errors.New("tried to write to a closed session")
	}

	s.conn.SetWriteDeadline(time.Now().Add(s.asura.Config.WriteWait))
	err := s.conn.WriteMessage(message.t, message.msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *Socket) closed() bool {
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()

	return !s.open
}

func (s *Socket) close() {
	if !s.closed() {
		s.rwmutex.Lock()
		s.open = false
		s.conn.Close()
		close(s.output)
		s.rwmutex.Unlock()
	}
}

func (s *Socket) ping() {
	s.writeRaw(&envelope{t: websocket.PingMessage, msg: []byte{}})
}

func (s *Socket) writePump() {
	ticker := time.NewTicker(s.asura.Config.PingPeriod)
	defer ticker.Stop()

loop:
	for {
		select {
		case msg, ok := <-s.output:
			if !ok {
				break loop
			}

			err := s.writeRaw(msg)
			if err != nil {
				s.asura.errorHandler(s, err)
				break loop
			}

			if msg.t == websocket.CloseMessage {
				break loop
			}

			if msg.t == websocket.TextMessage {
				s.asura.messageSentHandler(s, msg.msg)
			}

			if msg.t == websocket.BinaryMessage {
				s.asura.messageSentHandlerBinary(s, msg.msg)
			}
		case <-ticker.C:
			s.ping()
		}
	}
}

func (s *Socket) readPump() {
	s.conn.SetReadLimit(s.asura.Config.MaxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(s.asura.Config.PongWait))

	s.conn.SetPongHandler(func(string) error {
		s.conn.SetReadDeadline(time.Now().Add(s.asura.Config.PongWait))
		s.asura.pongHandler(s)
		return nil
	})

	if s.asura.closeHandler != nil {
		s.conn.SetCloseHandler(func(code int, text string) error {
			return s.asura.closeHandler(s, code, text)
		})
	}

	for {
		t, message, err := s.conn.ReadMessage()
		if err != nil {
			s.asura.errorHandler(s, err)
			break
		}

		if t == websocket.TextMessage {
			s.asura.messageHandler(s, message)
		}

		if t == websocket.BinaryMessage {
			s.asura.messageHandlerBinary(s, message)
		}
	}
}

// Write writes message to session.
func (s *Socket) EmitText(msg []byte) error {
	if s.closed() {
		return errors.New("session is closed")
	}

	s.writeMessage(&envelope{t: websocket.TextMessage, msg: msg})

	return nil
}

// WriteBinary writes a binary message to session.
func (s *Socket) EmitBinary(msg []byte) error {
	if s.closed() {
		return errors.New("session is closed")
	}

	s.writeMessage(&envelope{t: websocket.BinaryMessage, msg: msg})

	return nil
}

// WriteBinary writes a binary message to session.
func (s *Socket) Emit(msg []byte) error {
	return s.EmitBinary(msg)
}

func (s *Socket) Pong(msg []byte) error {
	if s.closed() {
		return errors.New("session is closed")
	}
	s.writeMessage(&envelope{t: websocket.PongMessage, msg: msg})
	return nil
}

// Close closes session.
func (s *Socket) Close() error {
	if s.closed() {
		return errors.New("session is already closed")
	}

	s.writeMessage(&envelope{t: websocket.CloseMessage, msg: []byte{}})

	return nil
}

// CloseWithMsg closes the session with the provided payload.
// Use the FormatCloseMessage function to format a proper close message payload.
func (s *Socket) CloseWithMsg(msg []byte) error {
	if s.closed() {
		return errors.New("session is already closed")
	}

	s.writeMessage(&envelope{t: websocket.CloseMessage, msg: msg})

	return nil
}

// Set is used to store a new key/value pair exclusivelly for this session.
// It also lazy initializes s.Keys if it was not used previously.
func (s *Socket) Set(key string, value interface{}) {
	if s.Keys == nil {
		s.Keys = make(map[string]interface{})
	}

	s.Keys[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (s *Socket) Get(key string) (value interface{}, exists bool) {
	if s.Keys != nil {
		value, exists = s.Keys[key]
	}

	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (s *Socket) MustGet(key string) interface{} {
	if value, exists := s.Get(key); exists {
		return value
	}

	panic("Key \"" + key + "\" does not exist")
}

// IsClosed returns the status of the connection.
func (s *Socket) IsClosed() bool {
	return s.closed()
}
