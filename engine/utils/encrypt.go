package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
)

func Md5sum(str []byte) string {
	l := md5.New()
	l.Write(str)
	return hex.EncodeToString(l.Sum(nil))
}

const S_APP_SECRET = "apsdfAJOJ(#@&($0809283JLJOOJ"

func SignEncrypt(plain string) string {
	var strPlain = S_APP_SECRET + base64.StdEncoding.EncodeToString([]byte(plain))
	return Md5sum([]byte(strPlain))
}

func CheckSign(sign, nonce string, mid uint32, seq int) bool {
	plain := fmt.Sprintf("%v%d%s", mid, seq, nonce)
	serverSign := SignEncrypt(plain)
	log.Println("CheckSign ", sign, serverSign)
	if sign == serverSign {
		return true
	}
	return false
}
