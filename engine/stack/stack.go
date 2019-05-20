package stack

import (
	"runtime"
	"strings"

	"nuvem/engine/logger"
)

func TraceStack() string {
	buf := make([]byte, 10000)
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	s := string(buf)

	// skip nano frames lines
	const skip = 7
	count := 0
	index := strings.IndexFunc(s, func(c rune) bool {
		if c != '\n' {
			return false
		}
		count++
		return count == skip
	})
	return s[index+1:]
}

// CatchPanic calls a function and returns the error if function paniced
func CatchPanic(f func()) (err interface{}) {
	defer func() {
		err = recover()
		if err != nil {
			logger.TraceError("%s panic: %s", f, err)
			//gwlog.TraceError("%s panic: %s", f, err)
		}
	}()

	f()
	return
}

// RunPanicless calls a function panic-freely
func RunPanicless(f func()) (panicless bool) {
	defer func() {
		err := recover()
		panicless = err == nil
		if err != nil {
			logger.TraceError("%s panic: %s", f, err)
			//gwlog.TraceError("%s panic: %s", f, err)
		}
	}()

	f()
	return
}

// RepeatUntilPanicless runs the function repeatly until there is no panic
func RepeatUntilPanicless(f func()) {
	for !RunPanicless(f) {
	}
}

// NextLargerKey finds the next key that is larger than the specified key,
// but smaller than any other keys that is larger than the specified key
func NextLargerKey(key string) string {
	return key + "\x00" // the next string that is larger than key, but smaller than any other keys > key
}
