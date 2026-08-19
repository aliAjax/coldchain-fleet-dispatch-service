package platform

import "time"

type Clock interface {
	Now() time.Time
}

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

func NewSystemClock() Clock { return systemClock{} }

type FrozenClock struct {
	At time.Time
}

func (f FrozenClock) Now() time.Time { return f.At }
