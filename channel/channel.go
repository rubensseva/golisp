package channel

import (
	"fmt"
	"sync"
)

type PeekableChannel struct {
	mutex  sync.Mutex
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

	fmt.Println("before", p.buffer)
	for _, el := range p.buffer {
		fmt.Println(el)
	}
	p.buffer = p.buffer[1:]
	fmt.Println("after", p.buffer)
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
