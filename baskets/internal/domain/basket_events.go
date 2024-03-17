package domain

type BasketStarted struct{
	Basket *Basket
}

func (BasketStarted) EventName() string{ return "baskets.BasketStarted"}

type BasketItemAdded struct{
	Basket *Basket
	Item Item
}

func (BasketItemAdded) EventName() string { return "baskets.BasketItemAdded"}

type BasketItemRemoved struct{
	Basket *Basket
	Item Item
}

func (BasketItemRemoved) EventName() string { return "baskets.BasketItemRemoved" }

type BasketCanceled struct{
	Basket *Basket
}

func (BasketCanceled) EventName() string { return "baskets.BasketCancelled" }

type BasketCheckout struct{
	Basket *Basket
}

func (BasketCheckout) EventName() string { return "baskets.BasketCheckout" }