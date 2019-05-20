package gate

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"nuvem/engine/asura"
	"nuvem/engine/coder"
	"nuvem/engine/logger"
	"nuvem/engine/proto"
	"nuvem/engine/tcp"
	"nuvem/engine/utils"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultAsyncMsgLen = 8192

type GateConfig struct {
	GameIdentify string
	ConnAddr     string
	Name         string
	SvrType      string
	Host         string
	PostUrl      string
}

type Gate struct {
	config    *GateConfig
	agent     tcp.Agent
	stopChan  chan bool
	conn      *tcp.TCPClient
	clientHub *ClientHub
	api       *sapi
}

var (
	_gate *Gate
)

func NewGate(config *GateConfig) {
	_gate = &Gate{
		config:    config,
		clientHub: NewClientHub(),
		stopChan:  make(chan bool),
		api:       newSAPI(),
	}
	_gate.run()
}

func GetGate() *Gate {
	return _gate
}

func (self *Gate) run() {
	self.conn = new(tcp.TCPClient)
	self.conn.Addr = self.config.ConnAddr
	self.conn.ConnNum = 1
	self.conn.ConnectInterval = 3 * time.Second
	self.conn.PendingWriteNum = defaultAsyncMsgLen
	self.conn.AutoReconnect = true
	self.conn.LenMsgLen = 2
	self.conn.MaxMsgLen = defaultAsyncMsgLen // math.MaxUint32
	self.conn.LittleEndian = false
	self.conn.NewAgent = func(conn *tcp.TCPConn) tcp.Agent {
		self.agent = &Agent{Conn: conn}
		return self.agent
	}

	if self.conn == nil {
		panic("runTCPClient err")
	}

	self.conn.Start()
}

func (self *Gate) StartTimer() {
	logger.Info("Gate:StartTimer at", time.Now().Unix())
	ticker := time.NewTicker(5 * time.Minute)
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				count := self.clientHub.ClientCount()
				self.RegisteGate(count)
			case stop := <-self.stopChan:
				if stop {
					logger.Info("Gate:StopTimer at", time.Now().Unix())
					return
				}
			}
		}
	}(ticker)
}

func (self *Gate) WriteMessage(mtype, sid string, msg []byte) error {
	var err error
	if self.agent != nil {
		err = self.agent.WriteMsg(mtype, sid, msg)
	} else {
		logger.Fatal("writeMessage no agent found")
	}

	if err != nil {
		logger.Error("WriteMessage err ", err)
	}

	return err
}

func (self *Gate) WriteAPI(sid string, msg []byte, c *ApiContext) error {
	if self.agent == nil {
		logger.Fatal("WriteAPI no agent found")
		return errors.New("no connection")
	}

	err := self.agent.WriteMsg(proto.MsgTypeAPI, sid, msg)
	if err != nil {
		logger.Error("WriteAPI err ", err)
		return err
	}
	if c != nil {
		self.api.push(sid, c)
	}

	return nil
}

//unionid orderid money
func (self *Gate) SendAPI(sid string, msg []byte) {
	apiContext := self.api.pop(sid)
	if apiContext == nil {
		logger.Error("SendAPI no context found", sid)
		return
	}
	defer close(apiContext.Quit)

	logger.Info("SendAPI", sid, string(msg))
	jsonmsg := make(coder.JSON)
	err := coder.ToJSON(msg, jsonmsg)
	if err != nil {
		logger.Error("SendAPI error", sid, err)
		apiContext.Ctext.JSON(http.StatusOK, gin.H{"error_code": http.StatusGatewayTimeout})
		return
	}

	outdata, ok := jsonmsg["data"].(map[string]interface{})
	if !ok {
		utils.DumpSocketData(outdata)
		apiContext.Ctext.JSON(http.StatusOK, gin.H{"error_code": http.StatusServiceUnavailable})
		return
	}
	logger.Info("SendAPI", sid, outdata)
	exchangeid := utils.GetInterfaceString("order", outdata)
	money := utils.GetInterfaceUint32("money", outdata)
	apiContext.Ctext.JSON(http.StatusOK, gin.H{"error_code": 0, "exchange_id": exchangeid, "exchange_num": money})
}

func (self *Gate) Destroy() {
	self.stopChan <- true
	close(self.stopChan)
	self.RegisteGate(99999)
	if self.conn != nil {
		self.conn.Close()
	}
}

func (self *Gate) RegisteGate(count int) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.PostForm(self.config.PostUrl,
		url.Values{
			"gameId": {self.config.GameIdentify},
			"name":   {self.config.Name},
			"value":  {self.config.Host},
			"count":  {utils.InterfaceToString(count)},
			"svr":    {self.config.SvrType},
		})
	if err != nil {
		logger.Error("RegisteGate PostForm err", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("RegisteGate ReadAll err", err)
	}
	logger.Debug("RegisteGate:body", string(body))
}

func (self *Gate) OnConnect(ws *asura.Socket) {
	self.clientHub.OnConnect(ws)
}

func (self *Gate) OnClose(ws *asura.Socket) {
	self.clientHub.OnClose(ws)
}

func (self *Gate) OnMessage(ws *asura.Socket, msg []byte) {
	self.clientHub.OnMessage(ws, msg)
}

func (self *Gate) CallBackMessage(mtype string, sid string, msg []byte) {
	self.clientHub.HandleMessage(mtype, sid, msg)
}
