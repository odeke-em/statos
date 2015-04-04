package statos

import (
	"io"
	"syscall"
)

// ReaderStatos implements the Read() interface
type ReaderStatos struct {
	done     bool
	finished uint64
	iterator io.Reader
}

func NewReader(rd io.Reader) *ReaderStatos {
	return &ReaderStatos{
		finished: 0,
		iterator: rd,
	}
}

func (r *ReaderStatos) Read(p []byte) (n int, err error) {
	n, err = r.iterator.Read(p)
	if err != nil && err != syscall.EINTR {
		r.done = true
	} else if n >= 0 {
		r.finished += uint64(n)
	}
	return
}

func (r *ReaderStatos) Progress() (uint64, bool) {
	return r.finished, r.done
}
