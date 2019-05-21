package coder

import (
	"errors"
	"fmt"
	"nuvem/engine/logger"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type JSON map[string]interface{}

func UnixSec() int64 {
	return time.Now().Unix()
}

func Reply(msgcode, errcode uint32, data interface{}, seq int) JSON {
	if seq == 0 {
		return JSON{
			"mid":  msgcode,
			"code": errcode,
			"time": time.Now().Unix(),
			"data": data,
		}
	}
	return JSON{
		"mid":  msgcode,
		"code": errcode,
		"time": time.Now().Unix(),
		"seq":  seq,
		"data": data,
	}
}

func ToBytes(reply interface{}) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	replybyte, err := json.Marshal(reply)
	if err != nil {
		logger.Error("ToBytes", err)
		return []byte(err.Error()), err
	}

	// replystr := string(replybyte)
	// replyss := []rune(replystr)

	// rsp := []byte{}
	// for _, v := range replyss {
	// 	rsp = append(rsp, byte(v))
	// 	rsp = append(rsp, byte(v>>8))
	// }
	//logger.Debug("ToBytes 11111", rsp)

	return replybyte, nil
}

func fromCharCode(r []rune, n int) string {
	i := 0
	for _, s := range r {
		if i == n {
			return string(s)
		}
		i++
	}
	return ""
}

func ToJSON(data []byte, jsondata JSON) error {
	// uarr := []uint16{}
	// for i := 0; i < len(data); i += 2 {
	// 	v := uint16(data[i]) + uint16(uint(data[i+1])<<8)
	// 	uarr = append(uarr, v)
	// }

	// str := ""
	// for _, v := range uarr {
	// 	if v == 0 {
	// 		str += "\\'\\0\\'"
	// 	} else {
	// 		str += string(v)
	// 	}
	// }
	// var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// err := json.UnmarshalFromString(str, &jsondata)
	// if err != nil {
	// 	return fmt.Errorf("ToJSON err: [%v]", err)
	// }
	// decodedata, err := charmap.ISO8859_1.NewDecoder().Bytes(data)
	// if err != nil {
	// 	return fmt.Errorf("ToJSON ISO8859_1 err: [%v]", err)
	// }
	// logger.Debug("data", data)
	// logger.Debug("decodedata", decodedata)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(data, &jsondata)
	if err != nil {
		return fmt.Errorf("ToJSON err: [%v]", err)
	}

	return nil
}

func Unpack(data []byte) (string, string, string, error) {
	var json = jsoniter.ConfigFastest
	jsondata := make(JSON)
	err := json.Unmarshal(data, &jsondata)
	if err != nil {
		return "", "", "", err
	}

	mtype, ok := jsondata["mtype"].(string)
	if !ok {
		mtype = ""
	}

	sid, ok := jsondata["sid"].(string)
	if !ok {
		sid = ""
	}

	outdata, ok := jsondata["data"].(string)
	if !ok {
		return "", "", "", errors.New("unpack data error")
	}

	return mtype, sid, outdata, nil
}

func Pack(mtype string, sid string, msg string) ([]byte, error) {
	var json = jsoniter.ConfigFastest
	data, err := json.Marshal(JSON{"mtype": mtype, "sid": sid, "data": msg})
	if err != nil {
		return nil, err
	}
	return data, nil
}
