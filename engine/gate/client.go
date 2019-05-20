package gate

import (
	"nuvem/engine/asura"
	"nuvem/engine/coder"
	"nuvem/engine/logger"
	"nuvem/engine/proto"
	"sync"
)

type Client struct {
	uid string
	ws  *asura.Socket
}

type ClientHub struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

func NewClientHub() *ClientHub {
	return &ClientHub{
		clients: make(map[string]*Client),
	}
}

func (self *ClientHub) OnConnect(ws *asura.Socket) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.clients[ws.Sid()] = &Client{ws: ws}
	logger.Debug("ClientHub:OnConnect", ws.Sid())
}

func (self *ClientHub) OnClose(ws *asura.Socket) {
	self.mu.Lock()
	defer self.mu.Unlock()
	_, ok := self.clients[ws.Sid()]
	if ok {
		msg := coder.ToBytes(coder.JSON{"mid": 1002}) //share.MsgCodeDisconnect
		err := GetGate().WriteMessage(proto.MsgTypeSignle, ws.Sid(), msg)
		if err != nil {
			logger.Error("OnMessage", ws.Sid(), err)
		}
	} else {
		logger.Error("ClientHub:OnClose error ", ws.Sid())
	}
	delete(self.clients, ws.Sid())
	logger.Debug("ClientHub:OnClose", ws.Sid())
}

func (self *ClientHub) OnMessage(ws *asura.Socket, msg []byte) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	_, ok := self.clients[ws.Sid()]
	if !ok {
		logger.Fatal("OnMessage no client found")
		ws.Close()
		return
	}

	err := GetGate().WriteMessage(proto.MsgTypeSignle, ws.Sid(), msg)
	if err != nil {
		logger.Error("OnMessage", ws.Sid(), err)
	}
}

func (self *ClientHub) HandleMessage(mtype string, sid string, msg []byte) {
	self.mu.RLock()
	defer self.mu.RUnlock()

	//logger.Debug("HandleMessage", sid, msg)
	if mtype == proto.MsgTypeSignle {
		client, ok := self.clients[sid]
		if !ok {
			logger.Error("HandleMessage no client", sid, string(msg))
			return
		}
		client.ws.Emit(msg)
	} else if mtype == proto.MsgTypeBCast {
		for _, client := range self.clients {
			//logger.Debug("BCAST", client.ws.Sid(), string(msg))
			client.ws.Emit(msg)
		}
	} else if mtype == proto.MsgTypeAPI {
		GetGate().SendAPI(sid, msg)
	} else {
		logger.Fatal("HandleMessage unknow message type", mtype)
	}
	// if sid == "" {
	// 	for _, client := range self.clients {
	// 		//logger.Debug("BCAST", client.ws.Sid(), string(msg))
	// 		client.ws.Emit(msg)
	// 	}
	// } else {
	// 	client, ok := self.clients[sid]
	// 	if !ok {
	// 		logger.Error("HandleMessage no client", sid, string(msg))
	// 		return
	// 	}
	// 	client.ws.Emit(msg)
	// }
}

func (self *ClientHub) ClientCount() int {
	self.mu.RLock()
	count := len(self.clients)
	self.mu.RUnlock()
	return count
}
