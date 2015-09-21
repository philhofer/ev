package ev

import (
	_ "runtime"
)

//go:noescape
func semrelease(l *uint32)

//go:noescape
func semacquire(l *uint32)
