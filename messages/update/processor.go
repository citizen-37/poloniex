package update

import (
	"fmt"
)

type (
	Update struct {
		parser  *Parser
		handler HandlerFunc
	}
)

func newUpdateMessage(handler HandlerFunc, parser *Parser) Update {
	return Update{
		parser:  parser,
		handler: handler,
	}
}

func (m Update) IsApplicable(message []interface{}) bool {
	return m.parser.IsApplicable(message)
}

func (m Update) Process(message []interface{}) error {
	payload, err := m.parser.Parse(message)
	if err != nil {
		return fmt.Errorf("cannot parse message: %w", err)
	}

	return m.handler(payload)
}
