package event2

import (
	"testing"
	"time"
)

func TestNewSimpleEvent(t *testing.T) {
	e := NewEventSystem()
	s := NewSimpleEvent(5*time.Second, true)
	e.Register(s)
	e.Run()
}
