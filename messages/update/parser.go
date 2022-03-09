package update

import (
	"fmt"
	"log"
)

type (
	Parser struct {
		itemParsers []ItemParser
	}
)

func NewParser(itemParsers []ItemParser) *Parser {
	return &Parser{
		itemParsers: itemParsers,
	}
}

func (p *Parser) IsApplicable(message []interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("failed to check message applicability: %v, %s", r, fmt.Sprintf("%s", message))
		}
	}()

	if len(message) < 2 {
		return false
	}

	body := message[2].([]interface{})

	if len(body) < 1 {
		return false
	}

	record := body[0].([]interface{})

	for _, itemParser := range p.itemParsers {
		if itemParser.IsApplicable(record) {
			return true
		}
	}

	return false
}

func (p *Parser) Parse(message []interface{}) (*OrderBookUpdate, error) {
	var dto OrderBookUpdate

	dto.ChannelId = p.extractChannelId(message)
	dto.SequenceNumber = p.extractSequenceNumber(message)

	err := p.parseAndApplyItems(&dto, message[2].([]interface{}))
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func (p *Parser) extractChannelId(message []interface{}) int64 {
	return int64(message[0].(float64))
}

func (p *Parser) extractSequenceNumber(message []interface{}) int64 {
	return int64(message[1].(float64))
}

func (p *Parser) parseAndApplyItems(dto *OrderBookUpdate, items []interface{}) error {
	var data []interface{}

	for _, item := range items {
		data = item.([]interface{})

		for _, itemParser := range p.itemParsers {
			if !itemParser.IsApplicable(data) {
				continue
			}

			err := itemParser.ParseAndApply(dto, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
