package pool

import "sync"

const (
	DefaultMaxBytesCap = 1 << 10 // 1k
	DefaultNewBytesCap = 512
)

var gloableBytesPool *BytesPool

type BytesPool struct {
	p       *sync.Pool
	maxCap  int
	initCap int
}

func (p *BytesPool) Get() *[]byte {
	return p.p.Get().(*[]byte)
}

func (p *BytesPool) Put(bz *[]byte) {
	if cap(*bz) > p.maxCap {
		return
	}
	b := (*bz)[:0]
	p.p.Put(&b)
}
