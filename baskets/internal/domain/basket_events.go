package domain

const (
	BasketStartedEvent     = "baskets.BasketStarted"
	BasketItemAddedEvent   = "baskets.BasketItemAdded"
	BasketItemRemovedEvent = "baskets.BasketItemRemoved"
	BasketCanceledEvent    = "baskets.BasketCanceled"
	BasketCheckedOutEvent  = "baskets.BasketCheckedOut"
)

type BasketStarted struct {
	CustomerID string
}

func (BasketStarted) EventName() string { return "baskets.BasketStarted" }

type BasketItemAdded struct {
	Item Item
}

func (BasketItemAdded) EventName() string { return "baskets.BasketItemAdded" }

type BasketItemRemoved struct {
	ProductID string
	Quantity  int
}

func (BasketItemRemoved) EventName() string { return "baskets.BasketItemRemoved" }

type BasketCanceled struct {
}

func (BasketCanceled) EventName() string { return "baskets.BasketCancelled" }

type BasketCheckedOut struct {
	PaymentID  string
	CustomerID string
	Items      map[string]Item
}

func (BasketCheckedOut) EventName() string { return "baskets.BasketCheckout" }
