package dump

import "fmt"

type (
	Dump struct {
		handler HandlerFunc
		parser  *Parser
	}
)

func newDumpMessageProcessor(handler HandlerFunc, parser *Parser) Dump {
	return Dump{
		handler: handler,
		parser:  parser,
	}
}

func (p Dump) IsApplicable(message []interface{}) bool {
	return p.parser.isApplicable(message)
}

func (p Dump) Process(message []interface{}) error {
	payload, err := p.parser.Parse(message)
	if err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	return p.handler(payload)
}
