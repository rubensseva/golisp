package channel

import (
	"sync"
)

type PeekableChannel struct {
	mutex sync.Mutex
	buffer []string
}

func NewPeekableChannel(size int) *PeekableChannel {
	return &PeekableChannel{
		buffer: make([]string, 0, size),
	}
}

func (p *PeekableChannel) Send(value string) {
	p.buffer = append(p.buffer, value)
}

func (p *PeekableChannel) Receive() (val string, available bool) {
	var z string
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.buffer) == 0 {
		return z, false
	}

	res := p.buffer[0]
	p.buffer = p.buffer[1:]
	return res, true

}

func (p *PeekableChannel) Peek() (val string, available bool) {
	var z string
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.buffer) == 0 {
		return z, false
	}

	res := p.buffer[0]
	return res, true
}
