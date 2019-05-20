package bridge

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"nuvem/engine/logger"
	"nuvem/engine/utils"
	"path"
	"time"
)

const defaultAsyncMsgLen = 81920

type Bridge struct {
	stopChan chan bool

	config *ConnConfig
	wsSvr  *Server
}

func (self *Bridge) StartTimer() {
	logger.Info("Bridge:StartTimer at", time.Now().Unix())
	ticker := time.NewTicker(3 * time.Minute)
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				count := 0
				if wHandler.binder != nil {
					count = wHandler.binder.BindCount()
				}
				logger.Debug("Bridge:Tick.........", count)
				self.RegisteBridge(count)
			case stop := <-self.stopChan:
				if stop {
					logger.Info("Bridge:StopTimer at", time.Now().Unix())
					return
				}
			}
		}
	}(ticker)
}

func (self *Bridge) Destroy() {
	self.RegisteBridge(99999)
	self.stopChan <- true
	close(self.stopChan)
	self.wsSvr.Shutdown()
	wHandler.Stop()
	logger.Info("Bridge:Destroy ^.^ ^.^ ^.^ ^.^ ")
}

func (self *Bridge) RegisteBridge(count int) {
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
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
		logger.Error("RegisteBridge PostForm err", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("RegisteBridge ReadAll err", err)
	}
	logger.Debug("RegisteBridge:body", string(body))
}

func NewBridge(authPath string, config *ConnConfig) *Bridge {
	_bridge := &Bridge{
		stopChan: make(chan bool),
		config:   config,
	}
	connecter := newConnecter(config.ConnAddr, _bridge)
	connecter.Start()

	wHandler = newWebsocketHandler(connecter)
	_bridge.wsSvr = NewServer(config.ServerAddr)
	go func() {
		if err := _bridge.wsSvr.ListenAndServe(path.Join(authPath, "public.pem"), path.Join(authPath, "private.key"), wHandler); err != nil {
			logger.Error("NewBridge:ListenAndServe error :", err)
		}
	}()
	logger.Info("NewBridge:ListenAndServe ", config.ServerAddr)
	return _bridge
}
