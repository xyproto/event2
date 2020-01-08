package event2

type Event interface {
	Trigger() error
	Hour() int
	Minute() int
	JustOnce() bool
}
