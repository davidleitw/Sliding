package slidingwindow

import (
	"fmt"
	"sync"
)

type kv struct {
	key          string
	defaultValue int

	next *kv
}

type kvlist struct {
	head   *kv
	tail   *kv
	length int
}

func (list *kvlist) insert(key string, value int) {
	if list.head == nil && list.tail == nil {
		list.head = &kv{key: key, defaultValue: value, next: nil}
		list.tail = list.head
		list.length = 1
		return
	}

	list.tail.next = &kv{key: key, defaultValue: value, next: nil}
	list.tail = list.tail.next
	list.length++
}

func (list *kvlist) setDefault(key string, value int) {
	head := list.head
	for head != nil {
		if head.key == key {
			head.defaultValue = value
			return
		}
		head = head.next
	}

	list.tail.next = &kv{key: key, defaultValue: value, next: nil}
	list.tail = list.tail.next
}

type uploadFunc func(int, map[string]int)

func WrapUploadFunc(upl func(int, map[string]int)) uploadFunc {
	return uploadFunc(upl)
}

type window struct {
	// When new start time != startTime, it means window needed to be update.
	startTime int64

	// mutex for updating data in window.
	mu sync.Mutex

	// Upload function will triggered when new start time != startTime
	upload uploadFunc

	// Define meta data key.
	defaultMedaDataKvlist *kvlist
	// Store meta data, and use upload function to save it when window is updated.
	// metaData can use in some complex scenario than counter.
	metaData map[string]int

	// Counter can count some simple information.
	counter int
}

func NewWindow(start int64, upl uploadFunc) *window {
	w := &window{
		startTime:             start,
		mu:                    sync.Mutex{},
		defaultMedaDataKvlist: new(kvlist),
		metaData:              make(map[string]int),
		counter:               0,
	}

	if upl != nil {
		w.registerUploadFunction(upl)
	}
	return w
}

func (w *window) registerUploadFunction(upl uploadFunc) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.upload = upl
}

func (w *window) atomicCounterAdd(delta int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.counter += delta
}

func (w *window) atomicMetaDataAdd(key string, delta int) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if value, ok := w.metaData[key]; ok {
		w.metaData[key] = value + delta
	} else {
		w.metaData[key] = delta
	}
}

func (w *window) setDefaultMetaDataKeys(keys []string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, key := range keys {
		w.defaultMedaDataKvlist.insert(key, 0)
	}

	head := w.defaultMedaDataKvlist.head
	for head != nil {
		fmt.Println(head.key, head.defaultValue)
		head = head.next
	}

}

func (w *window) setMedaDataDefaultKv(key string, value int) {
	w.mu.Lock()
	defer w.mu.Unlock()

}

func (w *window) checkStartTime(start int64) bool {
	return w.startTime == start
}

func (w *window) Update(start int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Update window start, and call its upload function.
	if w.upload != nil {
		w.upload(w.counter, w.metaData)
	}
	w.startTime = start

	// Initial data inside the window.
	w.reset()
}

// new window with new counter and meta data.
func (w *window) reset() {
	w.counter = 0

	head := w.defaultMedaDataKvlist.head
	for head != nil {
		w.metaData[head.key] = head.defaultValue
		head = head.next
	}
}
