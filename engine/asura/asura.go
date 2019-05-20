package asura

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Config asura configuration struct.
type Config struct {
	WriteWait         time.Duration // Milliseconds until write times out.
	PongWait          time.Duration // Timeout for waiting on pong.
	PingPeriod        time.Duration // Milliseconds between pings.
	MaxMessageSize    int64         // Maximum size in bytes of a message.
	MessageBufferSize int           // The max amount of messages that can be in a sessions buffer before it starts dropping them.
}

func newConfig() *Config {
	return &Config{
		WriteWait:         10 * time.Second,
		PongWait:          60 * time.Second,
		PingPeriod:        (60 * time.Second * 9) / 10,
		MaxMessageSize:    8192 * 2,
		MessageBufferSize: 8192,
	}
}

// Close codes defined in RFC 6455, section 11.7.
// Duplicate of codes from gorilla/websocket for convenience.
const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

// Duplicate of codes from gorilla/websocket for convenience.
var validReceivedCloseCodes = map[int]bool{
	// see http://www.iana.org/assignments/websocket/websocket.xhtml#close-code-number

	CloseNormalClosure:           true,
	CloseGoingAway:               true,
	CloseProtocolError:           true,
	CloseUnsupportedData:         true,
	CloseNoStatusReceived:        false,
	CloseAbnormalClosure:         false,
	CloseInvalidFramePayloadData: true,
	ClosePolicyViolation:         true,
	CloseMessageTooBig:           true,
	CloseMandatoryExtension:      true,
	CloseInternalServerErr:       true,
	CloseServiceRestart:          true,
	CloseTryAgainLater:           true,
	CloseTLSHandshake:            false,
}

type handleMessageFunc func(*Socket, []byte)
type handleErrorFunc func(*Socket, error)
type handleCloseFunc func(*Socket, int, string) error
type handleSessionFunc func(*Socket)
type filterFunc func(*Socket) bool

// asura implements a websocket manager.
type Asura struct {
	Config                   *Config
	Upgrader                 *websocket.Upgrader
	messageHandler           handleMessageFunc
	messageHandlerBinary     handleMessageFunc
	messageSentHandler       handleMessageFunc
	messageSentHandlerBinary handleMessageFunc
	errorHandler             handleErrorFunc
	closeHandler             handleCloseFunc
	connectHandler           handleSessionFunc
	disconnectHandler        handleSessionFunc
	pongHandler              handleSessionFunc
	hub                      *hub
}

// New creates a new asura instance with default Upgrader and Config.
func New() *Asura {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  20480,
		WriteBufferSize: 20480,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	hub := newHub()

	go hub.run()

	return &Asura{
		Config:                   newConfig(),
		Upgrader:                 upgrader,
		messageHandler:           func(*Socket, []byte) {},
		messageHandlerBinary:     func(*Socket, []byte) {},
		messageSentHandler:       func(*Socket, []byte) {},
		messageSentHandlerBinary: func(*Socket, []byte) {},
		errorHandler:             func(*Socket, error) {},
		closeHandler:             nil,
		connectHandler:           func(*Socket) {},
		disconnectHandler:        func(*Socket) {},
		pongHandler:              func(*Socket) {},
		hub:                      hub,
	}
}

// HandleConnect fires fn when a session connects.
func (m *Asura) HandleConnect(fn func(*Socket)) {
	m.connectHandler = fn
}

// HandleDisconnect fires fn when a session disconnects.
func (m *Asura) HandleDisconnect(fn func(*Socket)) {
	m.disconnectHandler = fn
}

// HandlePong fires fn when a pong is received from a session.
func (m *Asura) HandlePong(fn func(*Socket)) {
	m.pongHandler = fn
}

// HandleMessage fires fn when a text message comes in.
func (m *Asura) HandleMessage(fn func(*Socket, []byte)) {
	m.messageHandler = fn
}

// HandleMessageBinary fires fn when a binary message comes in.
func (m *Asura) HandleMessageBinary(fn func(*Socket, []byte)) {
	m.messageHandlerBinary = fn
}

// HandleSentMessage fires fn when a text message is successfully sent.
func (m *Asura) HandleSentMessage(fn func(*Socket, []byte)) {
	m.messageSentHandler = fn
}

// HandleSentMessageBinary fires fn when a binary message is successfully sent.
func (m *Asura) HandleSentMessageBinary(fn func(*Socket, []byte)) {
	m.messageSentHandlerBinary = fn
}

// HandleError fires fn when a session has an error.
func (m *Asura) HandleError(fn func(*Socket, error)) {
	m.errorHandler = fn
}

// HandleClose sets the handler for close messages received from the session.
// The code argument to h is the received close code or CloseNoStatusReceived
// if the close message is empty. The default close handler sends a close frame
// back to the session.
//
// The application must read the connection to process close messages as
// described in the section on Control Frames above.
//
// The connection read methods return a CloseError when a close frame is
// received. Most applications should handle close messages as part of their
// normal error handling. Applications should only set a close handler when the
// application must perform some action before sending a close frame back to
// the session.
func (m *Asura) HandleClose(fn func(*Socket, int, string) error) {
	if fn != nil {
		m.closeHandler = fn
	}
}

// HandleRequest upgrades http requests to websocket connections and dispatches them to be handled by the asura instance.
func (m *Asura) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	return m.HandleRequestWithKeys(w, r, nil)
}

// HandleRequestWithKeys does the same as HandleRequest but populates session.Keys with keys.
func (m *Asura) HandleRequestWithKeys(w http.ResponseWriter, r *http.Request, keys map[string]interface{}) error {
	if m.hub.closed() {
		return errors.New("asura instance is closed")
	}

	conn, err := m.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	id := generateSidBytes(16)
	sock := &Socket{
		id:      b32enc.EncodeToString(id),
		Keys:    keys,
		conn:    conn,
		output:  make(chan *envelope, m.Config.MessageBufferSize),
		asura:   m,
		open:    true,
		rwmutex: &sync.RWMutex{},
	}

	m.hub.register <- sock
	m.connectHandler(sock)

	go sock.writePump()
	sock.readPump()

	if !m.hub.closed() {
		m.hub.unregister <- sock
	}

	sock.close()
	m.disconnectHandler(sock)
	return nil
}

// Broadcast broadcasts a text message to all sessions.
func (m *Asura) Broadcast(msg []byte) error {
	if m.hub.closed() {
		return errors.New("asura instance is closed")
	}

	message := &envelope{t: websocket.TextMessage, msg: msg}
	m.hub.broadcast <- message
	return nil
}

// BroadcastFilter broadcasts a text message to all sessions that fn returns true for.
func (m *Asura) BroadcastFilter(msg []byte, fn func(*Socket) bool) error {
	if m.hub.closed() {
		return errors.New("asura instance is closed")
	}

	message := &envelope{t: websocket.TextMessage, msg: msg, filter: fn}
	m.hub.broadcast <- message
	return nil
}

// BroadcastOthers broadcasts a text message to all sessions except session s.
func (m *Asura) BroadcastOthers(msg []byte, s *Socket) error {
	return m.BroadcastFilter(msg, func(q *Socket) bool {
		return s != q
	})
}

// BroadcastMultiple broadcasts a text message to multiple sessions given in the sessions slice.
func (m *Asura) BroadcastMultiple(msg []byte, sockets []*Socket) error {
	for _, so := range sockets {
		if writeErr := so.EmitText(msg); writeErr != nil {
			return writeErr
		}
	}
	return nil
}

// BroadcastBinary broadcasts a binary message to all sessions.
func (m *Asura) BroadcastBinary(msg []byte) error {
	if m.hub.closed() {
		return errors.New("asura instance is closed")
	}

	message := &envelope{t: websocket.BinaryMessage, msg: msg}
	m.hub.broadcast <- message
	return nil
}

// BroadcastBinaryFilter broadcasts a binary message to all sessions that fn returns true for.
func (m *Asura) BroadcastBinaryFilter(msg []byte, fn func(*Socket) bool) error {
	if m.hub.closed() {
		return errors.New("asura instance is closed")
	}

	message := &envelope{t: websocket.BinaryMessage, msg: msg, filter: fn}
	m.hub.broadcast <- message
	return nil
}

// BroadcastBinaryOthers broadcasts a binary message to all sessions except session s.
func (m *Asura) BroadcastBinaryOthers(msg []byte, s *Socket) error {
	return m.BroadcastBinaryFilter(msg, func(q *Socket) bool {
		return s != q
	})
}

// Close closes the asura instance and all connected sessions.
func (m *Asura) Close() error {
	if m.hub.closed() {
		return errors.New("asura instance is already closed")
	}

	m.hub.exit <- &envelope{t: websocket.CloseMessage, msg: []byte{}}
	return nil
}

// CloseWithMsg closes the Asura instance with the given close payload and all connected sessions.
// Use the FormatCloseMessage function to format a proper close message payload.
func (m *Asura) CloseWithMsg(msg []byte) error {
	if m.hub.closed() {
		return errors.New("Asura instance is already closed")
	}

	m.hub.exit <- &envelope{t: websocket.CloseMessage, msg: msg}
	return nil
}

// Len return the number of connected sessions.
func (m *Asura) Len() int {
	return m.hub.len()
}

// IsClosed returns the status of the Asura instance.
func (m *Asura) IsClosed() bool {
	return m.hub.closed()
}

// FormatCloseMessage formats closeCode and text as a WebSocket close message.
func FormatCloseMessage(closeCode int, text string) []byte {
	return websocket.FormatCloseMessage(closeCode, text)
}
