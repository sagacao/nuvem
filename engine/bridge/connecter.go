package bridge

import (
	"nuvem/engine/logger"
	"nuvem/engine/tcp"
	"time"
)

type Connecter struct {
	stopChan       chan bool
	conn           *tcp.TCPClient
	connecterAgent *ConnecterAgent
}

func newConnecter(svrAddr string, bridge *Bridge) *Connecter {
	c := &Connecter{
		conn: new(tcp.TCPClient),
	}

	c.conn = new(tcp.TCPClient)
	c.conn.Addr = svrAddr
	c.conn.ConnNum = 1
	c.conn.ConnectInterval = 3 * time.Second
	c.conn.PendingWriteNum = defaultAsyncMsgLen
	c.conn.AutoReconnect = true
	c.conn.LenMsgLen = 2
	c.conn.MaxMsgLen = defaultAsyncMsgLen // math.MaxUint32
	c.conn.LittleEndian = false
	c.conn.NewAgent = func(tconn *tcp.TCPConn) tcp.Agent {
		c.connecterAgent = &ConnecterAgent{Conn: tconn}
		c.connecterAgent.AfterReadFunc = func(messageType string, connID string, data []byte) {
			wHandler.Dispatch(messageType, connID, data)
		}
		c.connecterAgent.AfterConnectFunc = func() {
			bridge.RegisteBridge(0)
		}
		return c.connecterAgent
	}

	return c
}

func (c *Connecter) Start() {
	c.conn.Start()
}

func (c *Connecter) Stop() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Connecter) SendToServer(mType, connID string, message []byte) {
	var err error
	if c.connecterAgent != nil {
		err = c.connecterAgent.WriteMsg(mType, connID, message)
	} else {
		logger.Info("Connecter:SendToServer no agent.")
	}

	if err != nil {
		logger.Error("Connecter:SendToServer write error:", err)
	}
}
