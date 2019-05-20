package tcp

import (
	"net"
	"nuvem/engine/logger"
	"sync"
	"time"
)

type TCPClient struct {
	sync.Mutex
	Addr            string
	ConnNum         int
	ConnectInterval time.Duration
	PendingWriteNum int
	AutoReconnect   bool
	NewAgent        func(*TCPConn) Agent
	conns           ConnSet
	wg              sync.WaitGroup
	closeFlag       bool

	// msg parser
	LenMsgLen    int
	MinMsgLen    uint32
	MaxMsgLen    uint32
	LittleEndian bool
	msgParser    *MsgParser
}

func (client *TCPClient) Start() {
	client.init()

	for i := 0; i < client.ConnNum; i++ {
		client.wg.Add(1)
		go client.connect()
	}
}

func (client *TCPClient) init() {
	client.Lock()
	defer client.Unlock()

	if client.ConnNum <= 0 {
		client.ConnNum = 1
		logger.Info("invalid ConnNum, reset to %v", client.ConnNum)
	}
	if client.ConnectInterval <= 0 {
		client.ConnectInterval = 3 * time.Second
		logger.Info("invalid ConnectInterval, reset to", client.ConnectInterval)
	}
	if client.PendingWriteNum <= 0 {
		client.PendingWriteNum = 1000
		logger.Info("invalid PendingWriteNum, reset to", client.PendingWriteNum)
	}
	if client.NewAgent == nil {
		logger.Fatal("NewAgent must not be nil")
	}
	if client.conns != nil {
		logger.Fatal("client is running")
	}

	client.conns = make(ConnSet)
	client.closeFlag = false

	// msg parser
	msgParser := NewMsgParser()
	msgParser.SetMsgLen(client.LenMsgLen, client.MinMsgLen, client.MaxMsgLen)
	msgParser.SetByteOrder(client.LittleEndian)
	client.msgParser = msgParser
}

func (client *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil || client.closeFlag {
			return conn
		}

		logger.Info("connect to", client.Addr, " error: %v", err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *TCPClient) connect() {
	defer client.wg.Done()

reconnect:
	conn := client.dial()
	if conn == nil {
		return
	}

	client.Lock()
	if client.closeFlag {
		client.Unlock()
		conn.Close()
		return
	}

	tcpConn := newTCPConn(conn, client.PendingWriteNum, client.msgParser)
	agent := client.NewAgent(tcpConn)
	client.conns[conn] = agent
	client.Unlock()
	agent.OnConnect()
	agent.Run()

	// cleanup
	tcpConn.Close()
	client.Lock()
	delete(client.conns, conn)
	client.Unlock()
	agent.OnClose()

	if client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto reconnect
	}
}

func (client *TCPClient) Close() {
	client.Lock()
	client.closeFlag = true
	for conn := range client.conns {
		conn.Close()
	}
	client.conns = nil
	client.Unlock()

	client.wg.Wait()
}
