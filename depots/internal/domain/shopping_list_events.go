package domain

const (
	ShoppingListCreatedEvent   = "depots.ShoppingListCreated"
	ShoppingListCanceledEvent  = "depots.ShoppingListCanceled"
	ShoppingListAssignedEvent  = "depots.ShoppingListAssigned"
	ShoppingListCompletedEvent = "depots.ShoppingListCompleted"
)

type ShoppingListCreated struct {
	ShoppingList *ShoppingList
}

func (ShoppingListCreated) EventName() string { return "depot.ShoppingListCreated" }

func (ShoppingListCreated) Key() string { return ShoppingListCreatedEvent}

type ShoppingListCanceled struct {
	ShoppingList *ShoppingList
}

func (ShoppingListCanceled) EventName() string { return "depot.ShoppingListCanceled" }

func (ShoppingListCanceled) Key() string { return ShoppingListCanceledEvent}

type ShoppingListAssigned struct {
	ShoppingList *ShoppingList
	BotID        string
}

func (ShoppingListAssigned) EventName() string { return "depot.ShoppingListAssigned" }

func (ShoppingListAssigned) Key() string { return ShoppingListAssignedEvent }

type ShoppingListCompleted struct {
	ShoppingList *ShoppingList
}

func (ShoppingListCompleted) EventName() string { return "depot.ShoppingListCompleted" }

func (ShoppingListCompleted) Key() string { return ShoppingListCompletedEvent }
