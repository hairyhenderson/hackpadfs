package tar

import (
	"sync/atomic"
)

// bufferPool maintains a collection of byte buffers with a maximum size.
// Used to control upper-bound memory usage. It's safe for concurrent use.
type bufferPool struct {
	buffers chan *buffer
	count   atomic.Int64
	size    atomic.Uint64
}

type buffer struct {
	pool *bufferPool
	Data []byte
}

func newBufferPool(bufferSize, maxBuffers uint64) *bufferPool {
	if maxBuffers == 0 {
		maxBuffers = 1
	}
	p := &bufferPool{
		buffers: make(chan *buffer, maxBuffers),
	}
	p.size.Store(bufferSize)
	p.addBuffer() // start with 1 buffer, ready to go
	return p
}

func (p *bufferPool) addBuffer() {
	for {
		count := p.count.Load()
		if int(count) == cap(p.buffers) {
			return // already at max buffers, no-op
		}
		if p.count.CompareAndSwap(count, count+1) {
			break // successfully provisioned slot for new buffer
		}
	}
	buf := &buffer{
		Data: make([]byte, p.size.Load()),
		pool: p,
	}
	p.buffers <- buf
}

// Wait acquires and returns a buffer. Be sure to call buffer.Done() to return it to the pool.
func (p *bufferPool) Wait() *buffer {
	select {
	case buf := <-p.buffers:
		return buf
	default:
		p.addBuffer()
		// may not always get the new buffer, but looping could allocate more buffers far too quickly
		return <-p.buffers
	}
}

// Done returns this buffer to the pool
func (b *buffer) Done() {
	b.pool.buffers <- b
}
