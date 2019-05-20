package wsclient

import (
	"fmt"
	"nuvem/engine/coder"
	"nuvem/engine/logger"
	"nuvem/engine/utils"

	"github.com/gorilla/websocket"
)

type Wsclient struct {
	unionid string
	conn    *websocket.Conn
}

func New(addr string) *Wsclient {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		logger.Error("Dial ", addr, " error:[", err, "]")
		return nil
	}

	_ws := &Wsclient{
		conn: conn,
	}

	return _ws
}

func (ws *Wsclient) process() {
	for {
		messageType, message, err := ws.conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		//log.Println("porcess", message)
		if messageType == websocket.CloseMessage {
			logger.Debug("close..................")
		} else {
			ws.handler(message)
		}
	}
}

func (ws *Wsclient) send(jsondata coder.JSON) {
	msg := coder.ToBytes(jsondata)
	err := ws.conn.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		logger.Error("Wsclient send", err)
	}
}

func (ws *Wsclient) handler(msg []byte) {
	jsondata := coder.JSON{}
	err := coder.ToJSON(msg, jsondata)
	if err != nil {
		logger.Error("Wsclient handler", err)
		return
	}

	mid := utils.GetInterfaceUint32("mid", jsondata)
	logger.Debug("handler msg", mid, string(msg))

	// jsoniface, ok := jsondata["data"]
	// if !ok {
	// 	logger.Error("Wsclient handler on data found")
	// 	return
	// }

	// jsonobj, ok2 := jsoniface.(map[string]interface{})
	// if !ok2 {
	// 	utils.DumpSocketData(jsoniface)
	// 	logger.Error("Wsclient handler on data found")
	// 	return
	// }

}
