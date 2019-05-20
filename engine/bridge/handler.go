package bridge

import (
	"errors"
	"net/http"
	"nuvem/engine/logger"
	"nuvem/engine/proto"

	"github.com/gorilla/websocket"
)

var (
	wHandler *websocketHandler
)

// websocketHandler defines to handle websocket upgrade request.
type websocketHandler struct {
	// upgrader is used to upgrade request.
	upgrader *websocket.Upgrader

	// connecter is used to connect game.
	connecter *Connecter

	// binder stores relations about websocket connection and userID.
	binder *binder

	// calcUserIDFunc defines to calculate userID by token. The userID will
	// be equal to token if this function is nil.
	calcUserIDFunc func(token string) (userID string, ok bool)
}

// First try to upgrade connection to websocket. If success, connection will
// be kept until client send close message or server drop them.
func (wh *websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer wsConn.Close()

	// handle Websocket request
	sock := NewSocket(wsConn)
	//logger.Debug(" >>>>>>>>>>>>> connected", userID)
	// bind
	wh.binder.Bind(sock.GetID(), sock)
	sock.AfterReadFunc = func(messageType int, rmessage []byte) {
		if wh.connecter != nil {
			wh.connecter.SendToServer(proto.MsgTypeSignle, sock.GetID(), rmessage)
		}

		// userID := sock.GetID() //rm.Token
		// //logger.Debug(" >>>>>>>>>>>>> message", userID)
		// // bind
		// wh.binder.Bind(userID, sock)
	}
	sock.BeforeCloseFunc = func() {
		// unbind
		wh.binder.Unbind(sock)
		//logger.Debug(" <<<<<<<<<<<<<<< closed", sock.GetID())
	}

	sock.Listen()
	//disconnect
	wh.binder.Unbind(sock)
	//logger.Debug(" <<<<<<<<<<<<<<< disconnect", sock.GetID())
	if err := sock.Close(); err != nil {
		logger.Error(err)
	}
}

// closeConns unbind conns filtered by userID and event and close them.
// The userID can't be empty, but event can be empty. The event will be ignored
// if empty.
func (wh *websocketHandler) closeConn(connID string) (int, error) {
	sock, ok := wh.binder.FindConn(connID)
	if !ok {
		return 0, errors.New("closeConn no connection")
	}
	sock.Close()

	return 1, nil
}

func (wh *websocketHandler) Dispatch(mType, connID string, message []byte) {
	if mType == proto.MsgTypeSignle {
		so, ok := wh.binder.FindConn(connID)
		if !ok {
			logger.Debug("websocketHandler:Dispatch no conn found", connID)
		} else {
			_, err := so.Write(message)
			if err != nil {
				logger.Error("websocketHandler:Dispatch :[", connID, "] err:", err)
			}
		}
	} else if mType == proto.MsgTypeBCast {
		wh.binder.ForEach(func(so *Socket) {
			if so != nil {
				_, err := so.Write(message)
				if err != nil {
					logger.Error("websocketHandler:Dispatch err:", err)
				}
			}
		})
	} else if mType == proto.MsgTypeAPI {
		//Do API
	}
}

func (wh *websocketHandler) Stop() {
	if wh.connecter != nil {
		wh.connecter.Stop()
	}

	if wh.binder != nil {

	}
}

func newWebsocketHandler(c *Connecter) *websocketHandler {
	return &websocketHandler{
		upgrader:  defaultUpgrader,
		connecter: c,
		binder: &binder{
			connID2Sockets: make(map[string]*clientConn),
		},
	}
}
