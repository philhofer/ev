package ev

import (
	"io"
	"syscall"
)

// This file depends upon the following
// OS-specific routines:
//
// sockfd(fd int) *evfd // arm duplex fd
// rfd(fd int) *evfd    // arm read-only fd
// wfd(fd int) *evfd    // arm write-only fd
// delfd(fd int)        // disarm fd
// pipe() (r, w, err)
//

const fdmagic = 0x812fa0b2

func abort(cause string, erno int) {
	panic(cause + syscall.Errno(erno).Error())
}

type evfd struct {
	magic uint64 // magic bytes
	fd    int    // system fd
	rsema uint32 // read semaphore
	wsema uint32 // write semaphore
	atEOF bool
}

//go:noescape
func writev(fd int, iov *iovec, nev int) int

//go:noescape
func readv(fd int, iov *iovec, nev int) int

type iovec struct {
	base *byte
	size uintptr
}

func (f *evfd) Fd() int {
	return f.fd
}

func (f *evfd) Close() error {
	if f.magic == 0 {
		return nil
	}
	delfd(f.fd)
	err := syscall.Close(f.fd)
	f.fd = -1
	f.magic = 0
	return err
}

func (f *evfd) Write(p []byte) (int, error) {
	if f.magic == 0 {
		return 0, ErrClosed
	}
	var nn int
	for nn < len(p) {
		n, err := syscall.Write(f.fd, p[nn:])
		if err != nil {
			n = 0
			switch err {
			case syscall.EINTR:
				continue
			case syscall.EAGAIN:
				semacquire(&f.wsema)
				continue
			default:
				return nn, err
			}
		}
		if n == 0 {
			return nn, io.ErrClosedPipe
		}
		nn += n
	}
	return nn, nil
}

func (f *evfd) Read(p []byte) (int, error) {
	if f.magic == 0 {
		return 0, ErrClosed
	} else if f.atEOF {
		return 0, io.EOF
	}

try:
	n, err := syscall.Read(f.fd, p)
	if err != nil {
		n = 0
		switch err {
		case syscall.EINTR:
			goto try
		case syscall.EAGAIN:
			semacquire(&f.rsema)
			goto try
		default:
			return 0, err
		}
	}
	if n == 0 {
		f.atEOF = true
		return 0, io.EOF
	}
	return n, nil
}

/*
func (f *evfd) Writev(p [][]byte) (int, error) {
	if f.magic == 0 {
		return 0, ErrClosed
	}

	// write iovecs 32 at a time,
	// which should be 1 syscall
	// in most cases.
	var nn int
	for len(p) > 0 {
		var iovp [32]iovec
		var nev int
		t := p
		if len(t) > len(iovp) {
			t = t[:len(iovp)]
		}
		var total int
		for _, b := range t {
			if len(b) == 0 {
				continue
			}
			iovp[nev].base = &b[0]
			iovp[nev].size = uintptr(len(b))
			total += len(b)
			nev++
		}
		list := iovp[:nev]
		for len(list) > 0 {
			ok := writev(f.fd, &list[0], len(list))
			if ok < 0 {
				switch syscall.Errno(-ok) {
				case syscall.EAGAIN:
					semacquire(&f.wsema)
					continue
				case syscall.EINTR:
					continue
				default:
					return nn, syscall.Errno(-ok)
				}
			}
			// not a short write: continue the outer loop
			if ok == total {
				nn += total
				break
			}
			// short write: move the first
			// pointer forward in the list
			s := 0
			for ok >= int(list[s].size) {
				ok -= int(list[s].size)
				s++
				nn += int(list[s].size)
			}
			if ok > 0 {
				list[s].base = (*byte)(unsafe.Pointer(uintptr(list[s].base) + uintptr(ok)))
				list[s].size -= uintptr(ok)
				nn += int(ok)
			}
			list = list[s:]
		}
		p = p[len(t):]
	}
	return nn, nil
}*/
