package events

type Fetcher interface {
	Fetch(l int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	Message
	Document
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
