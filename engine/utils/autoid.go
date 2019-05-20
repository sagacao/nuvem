package utils

type AutoInc struct {
	start, step uint32
	queue       chan uint32
	running     bool
}

func NewAutoInc(start, step uint32) (ai *AutoInc) {
	ai = &AutoInc{
		start:   start,
		step:    step,
		running: true,
		queue:   make(chan uint32, 4),
	}
	go ai.process()
	return
}

func (ai *AutoInc) process() {
	defer func() { recover() }()
	for i := ai.start; ai.running; i = i + ai.step {
		ai.queue <- i
	}
}

func (ai *AutoInc) Id() uint32 {
	return <-ai.queue
}

func (ai *AutoInc) Close() {
	ai.running = false
	close(ai.queue)
}
