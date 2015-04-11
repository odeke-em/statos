package statos

import (
	"io"
	"syscall"
)

// WriterStatos implements the Write() interface
type WriterStatos struct {
	iterator io.WriteCloser
	commChan chan int
}

func NewWriter(w io.WriteCloser) *WriterStatos {
	return &WriterStatos{
		commChan: make(chan int),
		iterator: w,
	}
}

func (w *WriterStatos) Write(p []byte) (n int, err error) {
	n, err = w.iterator.Write(p)

	if err != nil && err != syscall.EINTR {
		close(w.commChan)
	} else if n >= 0 {
		w.commChan <- n
	}
	return
}

func (w *WriterStatos) ProgressChan() chan int {
	return w.commChan
}
