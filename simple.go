package event2

import (
	"log"
	"time"
)

// SimpleEvent is a simple event that will trigger a function at HH:MM
type SimpleEvent struct {
	hour   int
	minute int
	once   bool
	f      func() error
}

// NewSimpleEvent will create a simple event that can be triggered
func NewSimpleEvent(in time.Duration, once bool) *SimpleEvent {
	now := time.Now()
	when := now.Add(in)
	f := func() error {
		log.Println("DEFAULT SIMPLE EVENT TRIGGER FUNCTION")
		return nil
	}
	return &SimpleEvent{when.Hour(), when.Minute(), once, f}
}

// SetTriggerFunction can be used for replacing the trigger function for a single event
func (se *SimpleEvent) SetTriggerFunction(f func() error) {
	se.f = f
}

// Trigger will call the trigger function stored in the SimpleEvent struct
func (se *SimpleEvent) Trigger() error {
	// Call the function in f
	return se.f()
}

// Hour will return the hour number for when the event should trigger
func (se *SimpleEvent) Hour() int {
	return se.hour
}

// Minute will return the minute number for when the event should trigger
func (se *SimpleEvent) Minute() int {
	return se.minute
}

// JustOnce returns true if the event should only ever trigger once
func (se *SimpleEvent) JustOnce() bool {
	return se.once
}
