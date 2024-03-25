package domain

type StoreCreated struct {
	Store *Store
}

func (StoreCreated) EventName() string { return "stores.StoreCreted" }

type StoreParticipationEnabled struct {
	Store *Store
}

func (StoreParticipationEnabled) EventName() string { return "stores.StoresParticipationEnabled" }

type StoreParticipationDisabled struct {
	Store *Store
}

func (StoreParticipationDisabled) EventName() string { return "stores.StoresParticipationDisabled" }
