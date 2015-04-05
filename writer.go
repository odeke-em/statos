package statos

import (
	"io"
	"syscall"
)

// WriterStatos implements the Write() interface
type WriterStatos struct {
	iterator io.WriteCloser
	done     bool
	finished uint64
	// Monotomically increasing number to track the number of Writes
	curWriteV uint64
	// Track the previous
	prevWriteV uint64
	nlast      int
}

func NewWriter(w io.WriteCloser) *WriterStatos {
	return &WriterStatos{
		nlast:      0,
		curWriteV:  0,
		prevWriteV: 0,
		finished:   0,
		iterator:   w,
	}
}

func (w *WriterStatos) Write(p []byte) (n int, err error) {
	n, err = w.iterator.Write(p)

	w.prevWriteV = w.curWriteV
	w.nlast = n

	if err != nil && err != syscall.EINTR {
		w.done = true
	} else if n >= 0 {
		w.curWriteV += 1
		w.finished += uint64(n)
	}
	return
}

func (w *WriterStatos) Progress() (nlast int, finished uint64, fresh, done bool) {
	return w.nlast, w.finished, w.curWriteV > w.prevWriteV, w.done
}
