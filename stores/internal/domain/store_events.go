package domain

const (
	StoreCreatedEvent = "stores.StoreCreated"
	StoreParticipationEnabledEvent = "stores.StoreParticipationEnabled"
	StoreParticipationDisabledEvent = "stores.StoreParticipationDisabled"
	StoreRebrandedEvent = "stores.StoreRebranded"
)

type StoreCreated struct {
	Name string
	Location string
}

func (StoreCreated) EventName() string { return "stores.StoreCreted" }

func (StoreCreated) Key() string { return StoreCreatedEvent }

type StoreParticipationToggled struct {
	Participating bool
}

type StoreRebranded struct{
	Name string
}


func (StoreRebranded) Key() string { return StoreRebrandedEvent }
