package event2

import (
	"log"
	"time"
)

type SimpleEvent struct {
	hour   int
	minute int
	once   bool
}

func NewSimpleEvent(in time.Duration, once bool) *SimpleEvent {
	now := time.Now()
	when := now.Add(in)
	return &SimpleEvent{when.Hour(), when.Minute(), once}
}

func (se *SimpleEvent) Trigger() error {
	log.Println("SIMPLE EVENT :D")
	// Success
	return nil
}

func (se *SimpleEvent) Hour() int {
	return se.hour
}

func (se *SimpleEvent) Minute() int {
	return se.minute
}

func (se *SimpleEvent) JustOnce() bool {
	return se.once
}
