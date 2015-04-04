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
}

func NewWriter(w io.WriteCloser) *WriterStatos {
	return &WriterStatos{
		finished: 0,
		iterator: w,
	}
}

func (w *WriterStatos) Write(p []byte) (n int, err error) {
	n, err = w.iterator.Write(p)

	if err != nil && err != syscall.EINTR {
		w.done = true
	} else if n >= 0 {
		w.finished += uint64(n)
	}
	return
}

func (w *WriterStatos) Progress() (uint64, bool) {
	return w.finished, w.done
}

func (w *WriterStatos) Close() error {
	err := w.iterator.Close()
	if err == nil {
		w.done = true
	}
	return err
}
