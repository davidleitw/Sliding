package slidingwindow

import (
	"fmt"
	"sync"
	"time"
)

type SlidingWindows interface {
	Sync()
}

type Slw struct {
	mu           sync.Mutex
	windows      []*window
	windowSize   int64
	windowLength int64

	currentIndex      int
	defaultUploadFunc uploadFunc
}

func NewSlidingWindows(size, length int64, upl uploadFunc) *Slw {
	slw := &Slw{
		mu:                sync.Mutex{},
		windows:           make([]*window, length),
		windowSize:        size,
		windowLength:      length,
		currentIndex:      0,
		defaultUploadFunc: upl,
	}

	beginTime := time.Now().UnixMilli()
	for offset := 0; offset < int(slw.windowLength); offset++ {
		start := beginTime + int64(offset)*slw.windowSize
		slw.windows[offset] = NewWindow(start, slw.defaultUploadFunc)
	}
	return slw
}

func (slw *Slw) SetNewWindowSize(size int64) {
	slw.mu.Lock()
	slw.windowSize = size
	slw.mu.Unlock()
}

func (slw *Slw) Sync() *Slw {
	slw.mu.Lock()
	defer slw.mu.Unlock()
	current := time.Now().UnixMilli()

	index := slw.getCurrentIndex(current)
	start := slw.getCurrentStart(current)

	slw.currentIndex = index
	// First round, create a new window with upload function.
	for {
		if !slw.windows[index].checkStartTime(start) {
			slw.windows[index].Update(start)
			return slw
		} else {
			return slw
		}
	}
}

func (slw *Slw) AtomicWindowCounterAdd(delta int32) *Slw {
	slw.windows[slw.currentIndex].atomicCounterAdd(delta)
	return slw
}

func (slw *Slw) AtomicWindowMetaDataAdd(key string, delta int) *Slw {
	slw.windows[slw.currentIndex].atomicMetaDataAdd(key, delta)
	return slw
}

func (slw *Slw) SetWindowMetaDataDefault(key string, value int) *Slw {
	slw.windows[slw.currentIndex].setDefaultMedaData(key, value)
	return slw
}

func (slw *Slw) getCurrentIndex(current int64) int {
	return int((current / slw.windowSize) % slw.windowLength)
}

func (slw *Slw) getCurrentStart(current int64) int64 {
	return current - (current % slw.windowSize)
}

func (slw *Slw) PrintInfo() int {
	cnt := 0
	for idx := 0; idx < int(slw.windowLength); idx++ {
		if slw.windows[idx] != nil {
			fmt.Println(idx, slw.windows[idx].counter, slw.windows[idx].metaData)
			cnt += int(slw.windows[idx].counter)
		}
	}
	fmt.Println("cnt = ", cnt)
	return cnt
}
