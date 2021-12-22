package slidingwindow

import (
	"fmt"
	"sync"
	"time"
)

var cnt int = 0

type parseFunc func(int, map[string]int)

func WrapparseFunc(upl func(int, map[string]int)) parseFunc {
	return parseFunc(upl)
}

type WindowChunk struct {
	winIndex int64
	winStart int64

	winCounter  int
	winMetaData map[string]int
}

type Parser interface {
	parseWindowChunk(WindowChunk)
}

type parser struct {
	slw *Slw
	mu  sync.Mutex

	lastIndex int64
	lastStart int64

	buffer chan WindowChunk

	// pf parseFunc
	cnt int
}

func Newparser(slw *Slw, length int64) *parser {
	parser := &parser{
		slw:       slw,
		mu:        sync.Mutex{},
		lastIndex: 0,
		lastStart: 0,
		buffer:    make(chan WindowChunk, int(length)),
		cnt:       0,
	}
	go parser.async()
	return parser
}

func (p *parser) setlastUpdate(index, start int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastIndex = index
	p.lastStart = start
}

// if start != window[index].start
// send old window data in parseWindowChunk
func (p *parser) parseWindowChunk(chunk WindowChunk) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.lastIndex != chunk.winIndex {
		p.buffer <- p.slw.windows[p.lastIndex].wrapWindowChunk()
		p.cnt += p.slw.windows[p.lastIndex].counter
	}

	p.cnt += int(chunk.winCounter)
	p.buffer <- chunk
}

func (p *parser) async() {
	for {
		select {
		case chunk := <-p.buffer:
			fmt.Println("Index ", chunk.winIndex, "->", time.UnixMilli(chunk.winStart).Format(time.RFC1123Z), ": ", chunk.winCounter, ", ", chunk.winMetaData)
			fmt.Println()
			cnt += chunk.winCounter
		default:
			continue
		}
	}
}

func (p *parser) c() int {
	return p.cnt
}
