package event2

import (
	"fmt"
	"sync"
	"time"
)

// 1. Create a new time system, and select granularity

// EventSys represents an event system.
// * granularity is how long it should wait at each loop in the main loop
// * events is a map from Event to a margin of error on each side of the time
//   when the event should kick in
type EventSys struct {
	granularity        time.Duration
	events             map[Event]time.Duration
	coolOffGranularity time.Duration
}

type Event struct {
	name   string
	hour   int
	minute int
}

// CoolOff is a list of all events that should not be triggered just yet
var coolOff []string
var mut sync.Mutex

func (e *Event) Trigger() error {
	fmt.Println("TRIGGERING " + e.name)
	// Placing in CoolOff
	mut.Lock()
	coolOff = append(coolOff, e.name)
	mut.Unlock()
	return nil
}

func (esys *EventSys) coolOffLoop() {
	// Every N seconds, remove the first entry from the coolOff slice
	for {
		mut.Lock()
		if len(coolOff) > 0 {
			coolOff = coolOff[1:]
		}
		mut.Unlock()
		time.Sleep(esys.coolOffGranularity)
	}
}

func New(granularity time.Duration) *EventSys {
	events := make(map[Event]time.Duration)
	return &EventSys{granularity, events, time.Minute * 5}
}

// Run will run the event system endlessly, in the foreground
func (esys *EventSys) eventLoop() error {
	for {
		// Check if any events should kick in at this point in time +- error margin, in seconds
		now := time.Now()
		for event, _ := range esys.events {
			//before := now.Sub(margin)
			//after := now.Add(margin)
			if now.Hour() == event.hour && now.Minute() == event.minute {
				event.Trigger()
			}
		}
		time.Sleep(esys.granularity)
	}
}

// Start will start the event system and return
func (esys *EventSys) Start() {
	go esys.coolOffLoop()
	go esys.eventLoop()
}
