package slidingwindow

import (
	"errors"
	"sync"
)

type uploadFunc func()

type window struct {
	// When new start time != startTime, it means window needed to be update.
	startTime int64

	// mutex for updating data in window.
	mu sync.Mutex

	// Upload function will triggered when new start time != startTime
	upload uploadFunc

	// Store meta data, and use upload function to save it when window is updated.
	// metaData can use in some complex scenario.
	metaData map[string]int

	// Counter can count some simple information.
	counter int
}

func NewWindow(start int64, upl uploadFunc) *window {
	w := &window{
		startTime: start,
		mu:        sync.Mutex{},
		metaData:  make(map[string]int),
		counter:   0,
	}

	if upl != nil {
		_ = w.registerUploadFunction(upl)
	}
	return w
}

func (w *window) registerUploadFunction(upl uploadFunc) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.upload != nil {
		return errors.New("already register upload function")
	}

	w.upload = upl
	return nil
}

func (w *window) checkStartTime(start int64) bool {
	return w.startTime == start
}

func (w *window) Update(start int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Update window start, and call its upload function.
	if w.upload != nil {
		w.upload()
	}
	w.startTime = start

	// Initial data inside the window.
	w.reset()
}

func (w *window) reset() {
	w.counter = 0
	w.metaData = make(map[string]int)
}