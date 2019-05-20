package bridge

import (
	"context"
	"net/http"
	"nuvem/engine/logger"
	"time"

	"github.com/gorilla/websocket"
)

const (
	serverDefaultWSPath = "/ws"
)

var defaultUpgrader = &websocket.Upgrader{
	ReadBufferSize:  65535,
	WriteBufferSize: 65535,
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

// Server defines parameters for running websocket server.
type Server struct {
	// Address for server to listen on
	Addr string

	// Path for websocket request, default "/ws".
	WSPath string

	// Upgrader is for upgrade connection to websocket connection using
	// "github.com/gorilla/websocket".
	//
	// If Upgrader is nil, default upgrader will be used. Default upgrader is
	// set ReadBufferSize and WriteBufferSize to 1024, and CheckOrigin always
	// returns true.
	Upgrader *websocket.Upgrader
	svr      *http.Server
	wh       *websocketHandler
}

// ListenAndServe listens on the TCP network address and handle websocket
// request.
func (s *Server) ListenAndServe(rsaCert, rsaKey string, handler *websocketHandler) error {
	s.wh = handler
	http.Handle(s.WSPath, s.wh)
	s.svr = &http.Server{
		Addr: s.Addr,
	}

	if rsaCert != "" && rsaKey != "" {
		return s.svr.ListenAndServeTLS(rsaCert, rsaKey)
		//return http.ListenAndServeTLS(s.Addr, rsaCert, rsaKey, nil)
	}

	return s.svr.ListenAndServe()
}

func (s *Server) Shutdown() {
	//s.wh.closeConns()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := s.svr.Shutdown(ctx); err != nil {
		logger.Fatal("WebSocket Server Shutdown:", err)
	}
	logger.Info("WebSocket Server exiting")
}

// Drop find connections by userID and event, then close them. The userID can't
// be empty. The event is ignored if it's empty.
func (s *Server) Drop(userID string) (int, error) {
	return s.wh.closeConn(userID)
}

// NewServer creates a new Server.
func NewServer(addr string) *Server {
	return &Server{
		Addr:   addr,
		WSPath: serverDefaultWSPath,
	}
}
