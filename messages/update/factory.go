package update

func NewUpdateMessageHandler(handler HandlerFunc) Update {
	p := NewParser([]ItemParser{
		NewItemUpdateParser(),
		NewItemTransactionParser(),
	})

	return newUpdateMessage(handler, p)
}
