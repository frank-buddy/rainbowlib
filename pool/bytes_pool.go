package pool

import "sync"

const (
	DefaultMaxBytesCap = 16 << 10 // 16k
	DefaultNewBytesCap = 512
)

var globalBytesPool *BytesPool
var bytesPoolOnce sync.Once

type BytesPool struct {
	p       *sync.Pool
	maxCap  int
	initCap int
}

func NewBytesPool(initCap, maxCap int) *BytesPool {
	if initCap < 1 {
		initCap = DefaultNewBytesCap
	}
	if maxCap < initCap {
		maxCap = initCap
	}
	return &BytesPool{
		p: &sync.Pool{
			New: func() interface{} {
				bz := make([]byte, 0, initCap)
				return &bz
			},
		},
		maxCap:  maxCap,
		initCap: initCap,
	}
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

func initGlobalBytesPool() {
	bytesPoolOnce.Do(func() {
		globalBytesPool = NewBytesPool(DefaultNewBytesCap, DefaultMaxBytesCap)
	})
}

func BytesPoolPut(bz *[]byte) {
	initGlobalBytesPool()
	globalBytesPool.Put(bz)
}

func BytesPoolGet() *[]byte {
	initGlobalBytesPool()
	return globalBytesPool.Get()
}

func SetBytesPoolMaxCap(maxCap int) {
	globalBytesPool.maxCap = maxCap
}
