package event2

import (
	"log"
	"sync"
	"time"
)

// EventSys represents an event system.
// * granularity is how long it should wait at each loop in the main loop
// * events is a map from Event to a margin of error on each side of the time
//   when the event should kick in
// * coolOffGranularity is how long the system should wait per cool-off
//   loop iteration
type EventSys struct {
	events             []Event
	granularity        time.Duration
	coolOffGranularity time.Duration
}

// CoolOff is a list of all events that should not be triggered just yet
var (
	coolOff    []Event
	coolOffMut sync.Mutex
)

func (sys *EventSys) coolOffLoop() {
	// Every N seconds, remove the first entry from the coolOff slice
	for {
		coolOffMut.Lock()
		if len(coolOff) > 0 {
			// If the event should be ran just once, move it to the back of the queue
			if coolOff[0].JustOnce() {
				tmp := coolOff[0]
				coolOff = coolOff[1:]
				coolOff = append(coolOff, tmp)
			} else {
				coolOff = coolOff[1:]
			}
		}
		coolOffMut.Unlock()
		time.Sleep(sys.coolOffGranularity)
	}
}

// NewEventSystem creates a new event system, where events can be registered
// and the event loop can be run.
func NewEventSystem() *EventSys {
	events := make([]Event, 0)
	granularity := time.Minute * 1
	coolOffDuration := time.Minute * 5
	return &EventSys{events, granularity, coolOffDuration}
}

// Register will register an event with the event system. Not thread safe.
func (sys *EventSys) Register(event Event) {
	// TODO: Mutex?
	sys.events = append(sys.events, event)
}

// eventLoop will run the event system endlessly, in the foreground
func (sys *EventSys) eventLoop() error {
	for {
		// Check if any events should kick in at this point in time +- error margin, in seconds
		now := time.Now()
	NEXT_EVENT:
		for _, event := range sys.events {
			// If the event is in the coolOff slice, skip for now
			for _, coolOffEvent := range coolOff {
				if coolOffEvent == event {
					log.Println("Skipping event that is in the cool-off period")
					continue NEXT_EVENT
				}
			}
			if now.Hour() == event.Hour() && now.Minute() == event.Minute() {
				log.Printf("Trigger event at %2d:%2d\n", now.Hour(), now.Minute())
				if event.Trigger() != nil {
					log.Println("event failed")
				}
				// Placing in the CoolOff slice,
				// which is handled by the cooloff-system
				coolOffMut.Lock()
				coolOff = append(coolOff, event)
				coolOffMut.Unlock()
			}
		}
		time.Sleep(sys.granularity)
	}
}

// Run will start the event system in the background and immediately return
func (sys *EventSys) Run() {
	go sys.coolOffLoop()
	go sys.eventLoop()
}

// Reset will remove a given event from the cool-off queue
func Reset(event Event) {
	coolOffMut.Lock()
	s := coolOff
	var i int
	for i2, e := range coolOff {
		if e == event {
			i = i2
			break
		}
	}
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	coolOff = s[:len(s)-1]
	coolOffMut.Unlock()
}
