package dump

func NewDumpMessageHandler(handler HandlerFunc) Dump {
	return newDumpMessageProcessor(handler, newParser())
}
