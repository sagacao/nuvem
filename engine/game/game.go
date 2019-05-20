package game

import (
	"nuvem/engine/tcp"
)

const defaultAsyncMsgLen = 81920

type Game struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32

	// tcp
	tcpServer    *tcp.TCPServer
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool
}

type FuncMessageHandle func(msgType string, sid string, msg []byte, agent tcp.Agent)

func NewGame(fn FuncMessageHandle, listenAddr string) *Game {
	_game := &Game{
		MaxConnNum:      100,
		PendingWriteNum: defaultAsyncMsgLen,
		MaxMsgLen:       4096,
		TCPAddr:         listenAddr, //ServerConfig.Env.SERVER_ADDR,
		LenMsgLen:       2,
		LittleEndian:    false,
	}
	_game.tcpServer = new(tcp.TCPServer)
	_game.tcpServer.Addr = _game.TCPAddr
	_game.tcpServer.MaxConnNum = _game.MaxConnNum
	_game.tcpServer.PendingWriteNum = _game.PendingWriteNum
	_game.tcpServer.LenMsgLen = _game.LenMsgLen
	_game.tcpServer.MaxMsgLen = _game.MaxMsgLen
	_game.tcpServer.LittleEndian = _game.LittleEndian
	_game.tcpServer.NewAgent = func(conn *tcp.TCPConn) tcp.Agent {
		a := &Agent{
			Conn:      conn,
			OnMessage: fn,
		}
		return a
	}

	controller.SetSvr(_game.tcpServer)

	if _game.tcpServer != nil {
		_game.tcpServer.Start()
	}
	return _game
}

func (g *Game) OnDestroy() {
	if g.tcpServer != nil {
		g.tcpServer.Close()
	}
}

func (g *Game) GetTCPSvr() *tcp.TCPServer {
	return g.tcpServer
}
