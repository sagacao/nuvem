package xmlparse

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type XmlRoot struct {
	XMLName xml.Name `xml:"root"` //标签名称
}

type XMLConfig struct {
	filename       string
	lastModifyTime int64
	rwLock         sync.RWMutex
	data           *XmlRoot
	notifyList     []Notifyer
}

func (self *XMLConfig) AddNotifyer(n Notifyer) {
	self.notifyList = append(self.notifyList, n)
}

func (self *XMLConfig) parse() (*XmlRoot, error) {
	data, err := ioutil.ReadFile(self.filename)
	if err != nil {
		fmt.Println("initXml read error:", err)
		return nil, err
	}

	var result XmlRoot
	err = xml.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("initXml Unmarshal error:", err)
		return nil, err
	}
	fmt.Println("initXml Success:", result)
	return &result, nil
}

func (self *XMLConfig) reload() {
	ticker := time.NewTicker(time.Second * 10)
	for _ = range ticker.C {
		func() {
			file, err := os.Open(self.filename)
			if err != nil {
				fmt.Printf("open %s failed,err:%v\n", self.filename, err)
				return
			}
			defer file.Close()
			fileInfo, err := file.Stat()
			if err != nil {
				fmt.Printf("stat %s failed,err:%v\n", self.filename, err)
				return
			}
			curModifyTime := fileInfo.ModTime().Unix()
			fmt.Printf("%v --- %v\n", curModifyTime, self.lastModifyTime)
			if curModifyTime > self.lastModifyTime {
				m, err := self.parse()
				if err != nil {
					fmt.Println("parse failed,err:", err)
					return
				}
				self.rwLock.Lock()
				self.data = m
				self.rwLock.Unlock()
				for _, n := range self.notifyList {
					n.Callback(self)
				}
				self.lastModifyTime = curModifyTime
			}
		}()
	}
}
