package bridge

import (
	"net"
	"nuvem/engine/coder"
	"nuvem/engine/logger"
	"nuvem/engine/tcp"
)

type ConnecterAgent struct {
	Conn             tcp.Conn
	AfterReadFunc    func(messageType string, connID string, data []byte)
	AfterConnectFunc func()
}

func (a *ConnecterAgent) Run() {
	for {
		data, err := a.Conn.ReadMsg()
		if err != nil {
			logger.Error("read message:", err)
			break
		}
		//logger.Debug("ReadMsg:", data)
		mtype, sid, msg, err := coder.Unpack(data)
		if err != nil {
			logger.Error("ConnecterAgent: read message err:", err)
		} else {
			if a.AfterReadFunc != nil {
				a.AfterReadFunc(mtype, sid, []byte(msg))
			}
		}
	}
}

func (a *ConnecterAgent) OnClose() {
	logger.Info("Connecter closed from:", a.RemoteAddr())
}

func (a *ConnecterAgent) OnConnect() {
	logger.Info("Connecter connect to: ", a.RemoteAddr(), "Success")
	if a.AfterConnectFunc != nil {
		a.AfterConnectFunc()
	}
}

func (a *ConnecterAgent) WriteMsg(mtype string, sid string, msg []byte) error {
	data, err := coder.Pack(mtype, sid, string(msg))
	if err != nil {
		logger.Error("pack msg err", err)
		return err
	}
	//logger.Debug("WriteMsg", sid, len(data), data)
	err = a.Conn.WriteMsg(data)
	if err != nil {
		logger.Error("write message", string(data), " error:", err)
		return err
	}
	return nil
}

func (a *ConnecterAgent) LocalAddr() net.Addr {
	return a.Conn.LocalAddr()
}

func (a *ConnecterAgent) RemoteAddr() net.Addr {
	return a.Conn.RemoteAddr()
}

func (a *ConnecterAgent) Close() {
	a.Conn.Close()
}

func (a *ConnecterAgent) Destroy() {
	a.Conn.Destroy()
}
