package bridge

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// VERSION returns current nano version
var VERSION = "0.0.1"

type ConnConfig struct {
	ServerAddr   string
	GameIdentify string
	ConnAddr     string
	Name         string
	SvrType      string
	Host         string
	PostUrl      string
}

var (
	// app represents the current server process
	app = &struct {
		name    string    // current application name
		startAt time.Time // startup time
	}{}

	// env represents the environment of the current process, includes
	// work path and config path etc.
	env = &struct {
		wd          string                   // working path
		die         chan bool                // wait for end application
		heartbeat   time.Duration            // heartbeat internal
		checkOrigin func(*http.Request) bool // check origin when websocket enabled
		debug       bool                     // enable debug
		wsPath      string                   // WebSocket path(eg: ws://127.0.0.1/wsPath)
	}{}
)

// init default configs
func init() {
	// application initialize
	app.name = strings.TrimLeft(filepath.Base(os.Args[0]), "/")
	app.startAt = time.Now()

	// environment initialize
	if wd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		env.wd, _ = filepath.Abs(wd)
	}

	env.die = make(chan bool)
	env.heartbeat = 30 * time.Second
	env.debug = false
	env.checkOrigin = func(_ *http.Request) bool { return true }
}
