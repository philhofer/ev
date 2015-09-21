package ev

import (
	"syscall"
	"time"
	"unsafe"
)

// kevent_t
type kvt struct {
	ident  uint64
	filter int16
	flags  uint16
	fflags uint32
	data   int64
	udata  unsafe.Pointer
}

//go:noescape
func kqueue() int32

//go:noescape
func kevent(fd int, ch *kvt, nch int, ev *kvt, nev int, ts *timespec) int32

func cloexec(fd int)

//go:noescape
func __pipe() (int, int)

type timespec struct {
	sec int64
	nsec int64
}

var (
	kq int
)

func init() {
	kq = int(kqueue())
	if kq < 0 {
		panic("kqueue: " + syscall.Errno(-kq).Error())
	}
	cloexec(kq)
	go run()
}

func sockfd(fd int) *evfd {
	e := &evfd{}
	ptr := unsafe.Pointer(e)
	kevf := uint16(syscall.EV_ADD | syscall.EV_CLEAR)
	var kvtp [2]kvt
	kvtp[0].ident = uint64(fd)
	kvtp[1].ident = uint64(fd)
	kvtp[0].flags = kevf
	kvtp[1].flags = kevf
	kvtp[0].filter = syscall.EVFILT_READ
	kvtp[1].filter = syscall.EVFILT_WRITE
	kvtp[0].udata = ptr 
	kvtp[1].udata = ptr

	e.fd = fd
	e.magic = fdmagic
	if ok := kevent(kq, &kvtp[0], 2, nil, 0, nil); ok < 0 {
		abort("kevent: ", -int(ok))
	}

	return e
}

func rfd(fd int) *evfd {
	e := &evfd{}
	ptr := unsafe.Pointer(e)
	kevf := uint16(syscall.EV_ADD | syscall.EV_CLEAR)
	var k kvt
	k.ident = uint64(fd)
	k.flags = kevf
	k.filter = syscall.EVFILT_READ
	k.udata = ptr

	e.fd = fd
	e.magic = fdmagic
	if ok := kevent(kq, &k, 1, nil, 0, nil); ok < 0 {
		abort("kevent: ", -int(ok))
	}

	return e
}

func wfd(fd int) *evfd {
	e := &evfd{}
	ptr := unsafe.Pointer(e)
	kevf := uint16(syscall.EV_ADD | syscall.EV_CLEAR)
	var k kvt
	k.ident = uint64(fd)
	k.flags = kevf
	k.filter = syscall.EVFILT_WRITE
	k.udata = ptr

	e.fd = fd
	e.magic = fdmagic
	if ok := kevent(kq, &k, 1, nil, 0, nil); ok < 0 {
		abort("kevent: ", -int(ok))
	}

	return e
}

func run() {
	for {
		var list [64]kvt
		ok := kevent(kq, nil, 0, &list[0], 64, nil)
		if ok < 0 {
			switch syscall.Errno(-ok) {
			case syscall.EAGAIN:
				println("kevent: EAGAIN")
				time.Sleep(5*time.Millisecond)
				continue
			case syscall.EINTR:
				continue
			default:
				abort("kevent: ", -int(ok))
			}
		}
		for i := range list[:ok] {
			k := &list[i]
			e := (*evfd)(k.udata)
			if e.magic != fdmagic {
				panic("bad magic")
			}
			switch k.filter {
			case syscall.EVFILT_READ:
				semrelease(&e.rsema)
			case syscall.EVFILT_WRITE:
				semrelease(&e.wsema)
			}
		}
	}
}

func delfd(fd int) {}


func pipe() (r int, w int, err syscall.Errno) {
	r, w = __pipe()
	if r < 0 {
		err, r, w = syscall.Errno(-r), 0, 0
		return
	}
	cloexec(r)
	cloexec(w)
	return
}


