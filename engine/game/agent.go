package game

import (
	"net"
	"nuvem/engine/coder"
	"nuvem/engine/logger"
	"nuvem/engine/tcp"
)

type Agent struct {
	Conn      tcp.Conn
	OnMessage FuncMessageHandle
}

func (a *Agent) Run() {
	for {
		data, err := a.Conn.ReadMsg()
		if err != nil {
			logger.Error("read message:", err)
			break
		}
		//logger.Debug("ReadMsg:", data)
		mtype, sid, msg, err := coder.Unpack(data)
		if err != nil {
			logger.Error("read message", err)
		} else {
			//logger.Debug("ReadMsg:", msg)
			if a.OnMessage != nil {
				a.OnMessage(mtype, sid, []byte(msg), a)
			}
			// a.Controller.HandleMassage(mtype, sid, []byte(msg), a)
		}
	}
}

func (a *Agent) OnClose() {
	logger.Error("Agent OnClose")
}

func (a *Agent) OnConnect() {
	logger.Debug("Game:Agent OnConnect")
}

func (a *Agent) WriteMsg(mtype, sid string, msg []byte) error {
	data, err := coder.Pack(mtype, sid, string(msg))
	if err != nil {
		logger.Error("pack msg err", err)
		return err
	}
	err = a.Conn.WriteMsg(data)
	if err != nil {
		logger.Error("write message", string(data), " error:", err)
		return err
	}
	return nil
}

func (a *Agent) LocalAddr() net.Addr {
	return a.Conn.LocalAddr()
}

func (a *Agent) RemoteAddr() net.Addr {
	return a.Conn.RemoteAddr()
}

func (a *Agent) Close() {
	a.Conn.Close()
}

func (a *Agent) Destroy() {
	a.Conn.Destroy()
}
