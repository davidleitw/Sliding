package slidingwindow

import (
	"sync"
)

type window struct {
	index int64

	// When new start time != startTime, it means window needed to be update.
	startTime int64

	// mutex for updating data in window.
	mu sync.Mutex

	// Define meta data key.
	defaultMetaList *Kv
	// Store meta data, and use upload function to save it when window is updated.
	// metaData can use in some complex scenario than counter.
	metaData map[string]int

	// Counter can count some simple information.
	counter int
}

func NewWindow(index, start int64) *window {
	w := &window{
		index:           index,
		startTime:       start,
		mu:              sync.Mutex{},
		defaultMetaList: new(Kv),
		metaData:        make(map[string]int),
		counter:         0,
	}
	return w
}

// Default all key with 0.
func (w *window) registerDefaultMetaKeys(keys []string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, key := range keys {
		w.defaultMetaList.setDefault(key, 0)
	}
}

func (w *window) setDefaultMetaKv(key string, value int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.defaultMetaList.setDefault(key, value)
}

func (w *window) removeDefaultKey(key string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.defaultMetaList.remove(key)
}

func (w *window) atomicCounterAdd(delta int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.counter += delta
}

func (w *window) atomicMetaDataAdd(key string, delta int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.metaData[key] += delta
}

func (w *window) checkStartTime(start int64) bool {
	return w.startTime == start
}

func (w *window) Update(start int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.startTime = start
	w.counter = 0

	kv := w.defaultMetaList.head
	for kv != nil {
		w.metaData[kv.key] = kv.defaultValue
		kv = kv.next
	}
}

func (w *window) wrapWindowChunk() WindowChunk {
	w.mu.Lock()
	defer w.mu.Unlock()

	chunk := WindowChunk{
		winIndex:   w.index,
		winStart:   w.startTime,
		winCounter: w.counter,
	}

	meta := make(map[string]int)

	for k, v := range w.metaData {
		meta[k] = v
	}
	chunk.winMetaData = meta
	return chunk
}
