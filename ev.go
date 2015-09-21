package ev

import (
	"errors"
)

var (
	ErrClosed = errors.New("ev: fd closed")
)

type Socket struct {
	*evfd
}
