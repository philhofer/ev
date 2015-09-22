package ev

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	epfd    int
	scratch [128]epollevent
)

type epollevent struct {
	events int32
	data   [8]byte // C and go align this field differently. (sigh)
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
		ok := epollwait(epfd, &scratch[0], len(scratch), -1)
		if ok < 0 {
			switch syscall.Errno(-ok) {
			case syscall.EAGAIN:
				println("epoll_wait: EAGAIN")
				time.Sleep(5 * time.Millisecond)
				continue
			case syscall.EINTR:
				continue
			default:
				abort("epoll_wait: ", -ok)
			}
		}
		for i := range scratch[:ok] {
			ev := &scratch[i]
			evfd := *(**evfd)(unsafe.Pointer(&ev.data))
			if evfd == nil {
				panic("nil evfd")
			}
			if evfd.magic != fdmagic {
				panic("bad magic")
			}
			if (ev.events & (syscall.EPOLLIN | syscall.EPOLLHUP | syscall.EPOLLRDHUP | syscall.EPOLLERR)) != 0 {
				semrelease(&evfd.rsema)
			}
			if (ev.events & (syscall.EPOLLOUT | syscall.EPOLLHUP | syscall.EPOLLERR)) != 0 {
				semrelease(&evfd.wsema)
			}
		}
	}
}

func sockfd(fd int) *evfd {
	e := &evfd{
		magic: fdmagic,
		fd:    fd,
	}
	ev := epollevent{
		events: syscall.EPOLLIN | syscall.EPOLLOUT | syscall.EPOLLET | syscall.EPOLLRDHUP,
	}
	*(**evfd)(unsafe.Pointer(&ev.data)) = e
	ok := epollctl(epfd, syscall.EPOLL_CTL_ADD, fd, &ev)
	if ok < 0 {
		abort("epoll_ctl: ", -ok)
	}
	return e
}

func rfd(fd int) *evfd {
	e := &evfd{
		magic: fdmagic,
		fd:    fd,
	}
	ev := epollevent{
		events: syscall.EPOLLIN | syscall.EPOLLET | syscall.EPOLLRDHUP,
	}
	*(**evfd)(unsafe.Pointer(&ev.data)) = e
	ok := epollctl(epfd, syscall.EPOLL_CTL_ADD, fd, &ev)
	if ok < 0 {
		abort("epoll_ctl: ", -ok)
	}
	return e
}

func wfd(fd int) *evfd {
	e := &evfd{
		magic: fdmagic,
		fd:    fd,
	}
	ev := epollevent{
		events: syscall.EPOLLOUT | syscall.EPOLLET,
	}
	*(**evfd)(unsafe.Pointer(&ev.data)) = e
	ok := epollctl(epfd, syscall.EPOLL_CTL_ADD, fd, &ev)
	if ok < 0 {
		abort("epoll_ctl: ", -ok)
	}
	return e
}

func delfd(fd int) {
	epollctl(epfd, syscall.EPOLL_CTL_DEL, fd, nil)
}

func pipe() (r int, w int, err syscall.Errno) {
	var dst [2]int32
	if ok := __pipe2(&dst, syscall.O_CLOEXEC|syscall.O_NONBLOCK); ok < 0 {
		err = syscall.Errno(-ok)
	} else {
		r, w = int(dst[0]), int(dst[1])
	}
	return
}

//go:noescape
func __pipe2(dst *[2]int32, flags int) int
