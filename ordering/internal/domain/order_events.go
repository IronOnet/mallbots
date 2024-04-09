package domain

const (
	OrderCreatedEvent = "ordering.OrderCreated"
	OrderCanceledEvent = "ordering.OrderCanceled"
	OrderReadiedEvent = "ordering.OrderReadied"
	OrderCompletedEvent = "ordering.OrderCompleted"
)

type OrderCreated struct{
	CustomerID string
	PaymentID string
	ShoppingID string
	Items []Item
}

func (OrderCreated) EventName() string { return "ordering.OrderCreated"}

func (OrderCreated) Key() string { return OrderCreatedEvent }

type OrderCanceled struct{
	CustomerID string
}

func (OrderCanceled) EventName() string { return "ordering.OrderCanceled"}

func (OrderCanceled) Key() string { return OrderCanceledEvent }


type OrderReadied struct{
	CustomerID string
	PaymentID string
	Total float64
}

func (OrderReadied) EventName() string { return "oredering.OrderReadied"}

func (OrderReadied) Key() string { return OrderReadiedEvent }

type OrderCompleted struct{
	InvoiceID string
}

func (OrderCompleted) EventName() string { return "ordering.OrderCompleted"}

func (OrderCompleted) Key() string { return OrderCompletedEvent }