package statos

import (
	"io"
	"syscall"
)

// ReaderStatos implements the Read() interface
type ReaderStatos struct {
	done bool
	// Monotomically increasing number to track the number of reads
	curReadV uint64
	// Track the previous
	prevReadV uint64
	lastRead  int
	finished  uint64
	iterator  io.Reader
}

func NewReader(rd io.Reader) *ReaderStatos {
	return &ReaderStatos{
		finished:  0,
		iterator:  rd,
		curReadV:  0,
		prevReadV: 0,
		lastRead:  0,
	}
}

func (r *ReaderStatos) Read(p []byte) (n int, err error) {
	n, err = r.iterator.Read(p)
	r.prevReadV = r.curReadV
	r.lastRead = n

	if err != nil && err != syscall.EINTR {
		r.done = true
	} else if n >= 0 {
		r.curReadV += 1
		r.finished += uint64(n)
	}
	return
}

func (r *ReaderStatos) Progress() (lastRead int, finished uint64, fresh, done bool) {
	return r.lastRead, r.finished, r.curReadV > r.prevReadV, r.done
}
