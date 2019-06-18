package utils

import (
	"errors"
	"math/rand"
	"time"
)

var (
	x0  uint32 = uint32(time.Now().UnixNano())
	a   uint32 = 1664525
	c   uint32 = 1013904223
	LCG chan uint32
)

const (
	PRERNG = 1024
)

// 全局快速随机数发生器，比标准库快，简单，可预生成
func init() {
	// LCG = make(chan uint32, PRERNG)
	// go func() {
	// 	for {
	// 		x0 = a*x0 + c
	// 		LCG <- x0
	// 	}

	// }()
	rand.Seed(time.Now().UnixNano())
}

func RandInterval(b1, b2 int32) int32 {
	if b1 == b2 {
		return b1
	}

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	return int32(rand.Int63n(max-min+1) + min)
}

func RandIntervalN(b1, b2 int32, n uint32) []int32 {
	if b1 == b2 {
		return []int32{b1}
	}

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	l := max - min + 1
	if int64(n) > l {
		n = uint32(l)
	}

	r := make([]int32, n)
	m := make(map[int32]int32)
	for i := uint32(0); i < n; i++ {
		v := int32(rand.Int63n(l) + min)

		if mv, ok := m[v]; ok {
			r[i] = mv
		} else {
			r[i] = v
		}

		lv := int32(l - 1 + min)
		if v != lv {
			if mv, ok := m[lv]; ok {
				m[v] = mv
			} else {
				m[v] = lv
			}
		}

		l--
	}

	return r
}

func Shift(sarray []interface{}) error {
	length := len(sarray)
	if length <= 0 {
		return errors.New("the length of the parameter sarray should not be less than 0")
	}

	for i := length - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		sarray[i], sarray[num] = sarray[num], sarray[i]
	}
	return nil
}

//GetRandomString
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
