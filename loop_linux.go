package ev

import (
	"os"
	"syscall"
	"time"
)

var (
	epfd    int
	scratch [128]epollevent
)

type epollevent struct {
	events uint32
	data   uint64
}

//go:noescape
func epollcreate1(flags int) int

//go:noescape
func epollwait(fd int, events *epollevent, nev int, ms int) int

//go:noescape
func epollctl(fd int, op int, nfd int, ev *epollevent) int

func init() {
	epfd = epollcreate1(syscall.EPOLL_CLOEXEC)
	if epfd < 0 {
		abort("epoll_create1: ", -epfd)
	}
	go run()
}

func run() {
	for {
		// block forever
		ok := epollwait(epfd, &scratch, len(scratch), 0)
		if ok < 0 {
			erno := syscall.Errno(-ok)
			switch erno {
			case sycall.EAGAIN:
				println("epoll_wait: EAGAIN")
				time.Sleep(5 * time.Millisecond)
				continue
			case syscall.EINTR:
				continue
			default:
				abort("epoll_wait: ", erno)
			}
		}
		for i := range scratch[:ok] {
			ev := &scratch[i]
			evfd := (*evfd)(unsafe.Pointer(uintptr(ev.data)))
			if evfd.magic != fdmagic {
				panic("bad magic")
			}
			if fd == nil {
				epollctl(l.epfd, EPOLL_CTL_DEL, ev.fd, nil)
				continue
			}
			if (ev.Events & syscall.EPOLLIN) != 0 {
				semrelease(&fd.rsema)
			}
			if (ev.Events & syscall.EPOLLOUT) != 0 {
				semrelease(&fd.wsema)
			}
			if (ev.Events & syscall.EPOLLERR) == syscall.EPOLLERR {
				semrelease(&fd.wsema)
				semrelease(&fd.rsema)
			}
		}
	}
}

func sockfd(fd int) *evfd {
	e := &evfd{}
	var ev epollevent
	ev.events = syscall.EPOLLIN | syscall.EPOLLOUT | syscall.EPOLLET
	ev.data = uint64(uintptr(unsafe.Pointer(e)))
	ok := epollctl(epfd, syscall.EPOLL_CTL_ADD, fd, &ev)
	if ok < 0 {
		abort("epoll_ctl: ", -ok)
	}
	e.magic = fdmagic
	e.fd = fd
	return e
}

func delfd(fd int) {
	epollctl(epfd, syscall.EPOLL_CTL_DEL, fd, nil)
}
