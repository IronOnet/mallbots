package domain

type BasketStatus string

const (
	BasketUnknown BasketStatus = ""
	BasketIsOpen BasketStatus = "open"
	BasketIsCanceled BasketStatus = "canceled"
	BasketIsCheckout BasketStatus = "checked_out"
)

func (s BasketStatus) String() string{
	switch s{
	case BasketIsOpen, BasketIsCanceled, BasketIsCheckout:
		return string(s)
	default:
		return ""
	}
}

func ToBasketStatus(status string) BasketStatus{
	switch status{
	case BasketIsOpen.String():
		return BasketIsOpen
	case BasketIsCanceled.String():
		return BasketIsCanceled
	case BasketIsCheckout.String():
		return BasketIsCheckout
	default:
		return BasketUnknown
	}
}