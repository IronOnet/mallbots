package domain

const (
	ProductAddedEvent = "stores.ProductAdded"
	ProductRebrandedEvent = "stores.ProductRebranded"
	ProductPriceIncreasedEvent = "stores.ProductPriceIncreased"
	ProductPriceDecreasedEvent = "stores.ProductPriceDecreased"
	ProductRemovedEvent = "stores.ProductRemoved"
)

type ProductAdded struct {
	StoreID string
	Name string
	Description string
	SKU string
	Price float64
}

func (ProductAdded) Key() string { return ProductAddedEvent }

func (ProductAdded) EventName() string { return "stores.ProductAdded" }

type ProductRebranded struct{
	Name string
	Description string
}

func (ProductRebranded) Key() string { return ProductRebrandedEvent }

type ProductPriceChanged struct{
	Delta float64
}

type ProductRemoved struct {
}

func (ProductRemoved) EventName() string { return "stores.ProductRemoved" }

func (ProductRemoved) Key() string { return ProductRemovedEvent }
