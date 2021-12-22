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
	mu              sync.Mutex
	windows         []*window
	windowSize      int64
	windowLength    int64
	windowRoundTime int64

	currentIndex int64
	parser       *parser
}

func NewSlidingWindows(size, length int64) *Slw {
	slw := &Slw{
		mu:              sync.Mutex{},
		windows:         make([]*window, length),
		windowSize:      size,
		windowLength:    length,
		windowRoundTime: size * length,
		currentIndex:    0,
	}

	slw.parser = Newparser(slw, slw.windowLength)
	beginTime := time.Now().UnixMilli()
	for offset := 0; offset < int(slw.windowLength); offset++ {
		start := beginTime + int64(offset)*slw.windowSize
		slw.windows[offset] = NewWindow(int64(offset), start)
		// slw.windows[offset] = NewWindow(int64(offset), 0)
	}
	return slw
}

func (slw *Slw) SetAllWindowsparserFunc(upl parseFunc) {
	slw.mu.Lock()
	defer slw.mu.Unlock()

	// for _, win := range slw.windows {
	// 	win.registerparserFunction(upl)
	// }
}

func (slw *Slw) RegisterDefaultMetaKeys(keys []string) *Slw {
	for _, win := range slw.windows {
		go win.registerDefaultMetaKeys(keys)
	}
	return slw
}

func (slw *Slw) SetDefaultMetaKv(key string, value int) *Slw {
	slw.mu.Lock()
	defer slw.mu.Unlock()

	for _, win := range slw.windows {
		win.setDefaultMetaKv(key, value)
	}
	return slw
}

func (slw *Slw) RemoveDefaultKey(key string) {
	slw.mu.Lock()
	defer slw.mu.Unlock()

	for _, win := range slw.windows {
		win.removeDefaultKey(key)
	}
}

func (slw *Slw) Sync() *Slw {
	slw.mu.Lock()
	defer slw.mu.Unlock()
	current := time.Now().UnixMilli()

	index := slw.getCurrentIndex(current)
	start := slw.getCurrentStart(current)

	slw.currentIndex = index
	// First round, create a new window with parser function.
	for {
		if !slw.windows[index].checkStartTime(start) {
			chunk := slw.windows[index].wrapWindowChunk()
			go func() {
				slw.parser.parseWindowChunk(chunk)
				slw.parser.setlastUpdate(index, start)
			}()

			slw.windows[index].Update(start)
			return slw
		} else {
			return slw
		}
	}
}

func (slw *Slw) AtomicWindowCounterAdd(delta int) *Slw {
	slw.windows[slw.currentIndex].atomicCounterAdd(delta)
	return slw
}

func (slw *Slw) AtomicWindowMetaDataAdd(key string, delta int) *Slw {
	slw.windows[slw.currentIndex].atomicMetaDataAdd(key, delta)
	return slw
}

func (slw *Slw) getCurrentIndex(current int64) int64 {
	return (current / slw.windowSize) % slw.windowLength
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
	fmt.Println("cnt = ", cnt, ", p.cnt = ", slw.parser.c())
	return cnt + slw.parser.c()
}

func (slw *Slw) Stop() {

}
