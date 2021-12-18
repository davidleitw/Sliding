package slidingwindow

import (
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
}

func NewSlidingWindows(size, length int64) *Slw {
	return &Slw{
		mu:           sync.Mutex{},
		windows:      make([]*window, length),
		windowSize:   size,
		windowLength: length,
	}
}

func (slw *Slw) SetNewWindowSize(size int64) {
	slw.mu.Lock()
	slw.windowSize = size
	slw.mu.Unlock()
}

func (slw *Slw) Sync() {
	current := time.Now().UnixMilli()

	index := slw.getCurrentIndex(current)
	start := slw.getCurrentStart(current)

	// First round, create a new window with upload function.
	if slw.windows[index] == nil {
		slw.mu.Lock()
		slw.windows[index] = NewWindow(start, nil)
		slw.mu.Unlock()
		return
	}

	// start != slw.windows[index].startTime
	if !slw.windows[index].checkStartTime(start) {
		slw.windows[index].Update(start)
		return
	}
}

func (slw *Slw) AtomicWindowCounterAdd(index int, delta int32) {
	slw.windows[index].atomicCounterAdd(delta)
}

func (slw *Slw) SetWindowMetaDataDefault(index int, key string, value int) {
	slw.windows[index].setDefaultMedaData(key, value)
}

func (slw *Slw) AtomicWindowMetaDataAdd(index int, key string, delta int) int {
	return slw.windows[index].atomicMetaDataAdd(key, delta)
}

func (slw *Slw) getCurrentIndex(current int64) int {
	return int((current / slw.windowSize) % slw.windowLength)
}

func (slw *Slw) getCurrentStart(current int64) int64 {
	return current - (current % slw.windowSize)
}
