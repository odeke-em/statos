package statos

import (
	"io"
	"syscall"
)

// WriteCloserStatos implements the Write() interface
type WriteCloserStatos struct {
	iterator io.WriteCloser
	done     bool
	finished uint64
}

func NewWriteCloser(w io.WriteCloser) *WriteCloserStatos {
	return &WriteCloserStatos{
		finished: 0,
		iterator: w,
	}
}

func (w *WriteCloserStatos) Write(p []byte) (n int, err error) {
	n, err = w.iterator.Write(p)

	if err != nil && err != syscall.EINTR {
		w.done = true
	} else if n >= 0 {
		w.finished += uint64(n)
	}
	return
}

func (w *WriteCloserStatos) Progress() (uint64, bool) {
	return w.finished, w.done
}

func (w *WriteCloserStatos) Close() error {
	err := w.iterator.Close()
	if err == nil {
		w.done = true
	}
	return err
}
