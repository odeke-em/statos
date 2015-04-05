package statos

import (
	"io"
	"syscall"
)

// WriteCloserStatos implements the Write() interface
type WriteCloserStatos struct {
	iterator io.WriteCloser
	done     bool
	// Monotomically increasing number to track the number of Writes
	curWriteV uint64
	// Track the previous
	prevWriteV uint64
	nlast      int
	finished   uint64
}

func NewWriteCloser(w io.WriteCloser) *WriteCloserStatos {
	return &WriteCloserStatos{
		nlast:      0,
		curWriteV:  0,
		prevWriteV: 0,
		finished:   0,
		iterator:   w,
	}
}

func (w *WriteCloserStatos) Write(p []byte) (n int, err error) {
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

func (w *WriteCloserStatos) Progress() (nlast int, finished uint64, fresh, done bool) {
	return w.nlast, w.finished, w.curWriteV > w.prevWriteV, w.done
}

func (w *WriteCloserStatos) Close() error {
	err := w.iterator.Close()
	if err == nil {
		w.done = true
	}
	return err
}
