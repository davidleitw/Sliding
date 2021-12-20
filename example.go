package main

import (
	"fmt"
	"sync"

	"github.com/davidleitw/Sliding/pkg/slidingwindow"
)

var cnt = 0
var wg sync.WaitGroup

func upload(counter int, metadata map[string]int) {
	cnt += int(counter)
	fmt.Println(counter, metadata)
}

func thread(slw *slidingwindow.Slw, tid int) {
	for i := 0; i < 40000000; i++ {
		slw.Sync().AtomicWindowCounterAdd(1).AtomicWindowMetaDataAdd(fmt.Sprintf("thread %d", tid), 2)
		// time.Sleep(3 * time.Microsecond)
	}
	defer wg.Done()
}

func main() {
	// Ten windows, each has 100 ms.
	slw := slidingwindow.NewSlidingWindows(250, 8, slidingwindow.WrapUploadFunc(upload))
	slw.SetDefaultMetaKv("thread 0", 0).SetDefaultMetaKv("thread 1", 0).SetDefaultMetaKv("thread 2", 0).SetDefaultMetaKv("thread 3", 0)
	fmt.Println(slw)

	wg.Add(4)
	go thread(slw, 0)
	go thread(slw, 1)
	go thread(slw, 2)
	go thread(slw, 3)
	wg.Wait()

	fmt.Println(cnt)
	wincnt := slw.PrintInfo()
	fmt.Println(cnt + wincnt)

	defer func(slw *slidingwindow.Slw) {
		if err := recover(); err != nil {
			fmt.Println("Panic catch!")
			slw.PrintInfo()
		}
	}(slw)
}
