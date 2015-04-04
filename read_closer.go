package statos

import (
	"io"
	"syscall"
)

// ReadCloserStatos implements the Read() interface
type ReadCloserStatos struct {
	done     bool
	finished uint64
	iterator io.ReadCloser
}

func NewReadCloser(rd io.ReadCloser) *ReadCloserStatos {
	return &ReadCloserStatos{
		finished: 0,
		iterator: rd,
	}
}

func (r *ReadCloserStatos) Read(p []byte) (n int, err error) {
	n, err = r.iterator.Read(p)
	if err != nil && err != syscall.EINTR {
		r.done = true
	} else if n >= 0 {
		r.finished += uint64(n)
	}
	return
}

func (r *ReadCloserStatos) Progress() (uint64, bool) {
	return r.finished, r.done
}

func (r *ReadCloserStatos) Close() error {
	err := r.iterator.Close()
	if err == nil {
		r.done = true
	}
	return err
}
