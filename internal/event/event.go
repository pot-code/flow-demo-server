package event

type Event interface {
	Topic() string
}

type EventBus interface {
	Publish(e Event)
}
