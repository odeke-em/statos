package statos

import (
	"fmt"
	"io"
	"syscall"
)

// ReadCloserStatos implements the Read() interface
type ReadCloserStatos struct {
	done bool
	// Monotomically increasing number to track the number of reads
	curReadV uint64
	// Track the previous
	prevReadV uint64
	finished  uint64
	iterator  io.ReadCloser
}

func NewReadCloser(rd io.ReadCloser) *ReadCloserStatos {
	return &ReadCloserStatos{
		curReadV:  0,
		prevReadV: 0,
		finished:  0,
		iterator:  rd,
	}
}

func (r *ReadCloserStatos) Read(p []byte) (n int, err error) {
	n, err = r.iterator.Read(p)

	r.prevReadV = r.curReadV

	fmt.Println(n, err)
	if err != nil && err != syscall.EINTR {
		r.done = true
	} else if n >= 0 {
		r.curReadV += 1
		r.finished += uint64(n)
	}
	return
}

func (r *ReadCloserStatos) Progress() (finished uint64, fresh, done bool) {
	return r.finished, r.curReadV > r.prevReadV, r.done
}

func (r *ReadCloserStatos) Close() error {
	err := r.iterator.Close()
	if err == nil {
		r.done = true
	}
	return err
}
