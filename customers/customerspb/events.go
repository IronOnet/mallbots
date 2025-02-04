package customerspb

import (
	"github.com/irononet/mallbots/internal/registry"
	"github.com/irononet/mallbots/internal/registry/serdes"
)

const (
	CustomerAggregateChannel = "mallbots.customers.events.Customer"

	CustomerRegisteredEvent = "customersapi.CustomerRegistered"
	CustomerSmsChangedEvent = "customersapi.CustomersSmsChanged"
	CustomerEnabledEvent = "customersapi.CustomerEnabled"
	CustomerDisabledEvent = "customersapi.CustomerDisabled"
)

func Registrations(reg registry.Registry) error{
	serde := serdes.NewProtoSerde(reg)

	// store events
	if err := serde.Register(&CustomerRegistered{}); err != nil{
		return err
	}

	if err := serde.Register(&CustomerSmsChanged{}); err != nil{
		return err
	}

	if err := serde.Register(&CustomerEnabled{}); err != nil{
		return err
	}
	
	if err := serde.Register(&CustomerDisabled{}); err != nil{
		return err
	}

	return nil
}

func  (*CustomerRegistered) Key() string { return CustomerRegisteredEvent }
func (*CustomerSmsChanged) Key() string { return CustomerSmsChangedEvent }
func (*CustomerEnabled) Key() string { return CustomerEnabledEvent }
func (*CustomerDisabled) Key() string { return CustomerDisabledEvent }