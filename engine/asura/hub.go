package asura

import (
	"sync"
)

type hub struct {
	sockets    map[*Socket]bool
	broadcast  chan *envelope
	register   chan *Socket
	unregister chan *Socket
	exit       chan *envelope
	open       bool
	rwmutex    *sync.RWMutex
}

func newHub() *hub {
	return &hub{
		sockets:    make(map[*Socket]bool),
		broadcast:  make(chan *envelope),
		register:   make(chan *Socket),
		unregister: make(chan *Socket),
		exit:       make(chan *envelope),
		open:       true,
		rwmutex:    &sync.RWMutex{},
	}
}

func (h *hub) run() {
loop:
	for {
		select {
		case s := <-h.register:
			h.rwmutex.Lock()
			h.sockets[s] = true
			h.rwmutex.Unlock()
		case s := <-h.unregister:
			if _, ok := h.sockets[s]; ok {
				h.rwmutex.Lock()
				delete(h.sockets, s)
				h.rwmutex.Unlock()
			}
		case m := <-h.broadcast:
			h.rwmutex.RLock()
			for s := range h.sockets {
				if m.filter != nil {
					if m.filter(s) {
						s.writeMessage(m)
					}
				} else {
					s.writeMessage(m)
				}
			}
			h.rwmutex.RUnlock()
		case m := <-h.exit:
			h.rwmutex.Lock()
			for s := range h.sockets {
				s.writeMessage(m)
				delete(h.sockets, s)
				s.Close()
			}
			h.open = false
			h.rwmutex.Unlock()
			break loop
		}
	}
}

func (h *hub) closed() bool {
	h.rwmutex.RLock()
	defer h.rwmutex.RUnlock()
	return !h.open
}

func (h *hub) len() int {
	h.rwmutex.RLock()
	defer h.rwmutex.RUnlock()

	return len(h.sockets)
}
