package ddd

const (
	AggregateNameKey    = "aggregate-name"
	AggregateIdKey      = "aggregate-id"
	AggregateVersionKey = "aggregate-version"
)

type (
	AggregateNamer interface {
		AggregateName() string
	}

	Eventer interface {
		AddEvent(string, EventPayload, ...EventOption)
		Events() []AggregateEvent
		ClearEvents()
	}

	Aggregate struct {
		Entity
		events []AggregateEvent
	}

	AggregateEvent interface {
		Event
		AggregateName() string
		AggregateID() string
		AggregateVersion() int
	}

	aggregateEvent struct {
		event
	}
)

var _ interface {
	AggregateNamer
	Eventer
} = (*Aggregate)(nil)

func NewAggregate(id, name string) Aggregate {
	return Aggregate{
		Entity: NewEntity(id, name),
		events: make([]AggregateEvent, 0),
	}
}

func (a Aggregate) AggregateName() string    { return a.name }
func (a Aggregate) Events() []AggregateEvent { return a.events }
func (a *Aggregate) ClearEvents()             { a.events = []AggregateEvent{} }

func (a *Aggregate) AddEvent(name string, payload EventPayload, options ...EventOption) {
	options = append(options,
	Metadata{
		AggregateNameKey: a.name,
		AggregateIdKey: a.id,
	},)

	a.events = append(a.events, 
		aggregateEvent{
			event: newEvent(name, payload, options...),
		},)
}

func (a *Aggregate) setEvents(events []AggregateEvent) { a.events = events }

func (a aggregateEvent) AggregateName() string { return a.metadata.Get(AggregateNameKey).(string)}
func (a aggregateEvent) AggregateID() string { return a.metadata.Get(AggregateIdKey).(string)}
func (a aggregateEvent) AggregateVersion() int { return a.metadata.Get(AggregateVersionKey).(int)}

type AggregateBase struct {
	ID     string
	events []Event
}

func (a *AggregateBase) GetID() string {
	return a.ID
}

func (a *AggregateBase) AddEvent(event Event) {
	a.events = append(a.events, event)
}

func (a AggregateBase) GetEvents() []Event {
	return a.events
}
