package dump

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type (
	Parser struct {
	}
)

func newParser() *Parser {
	return &Parser{}
}

func (p *Parser) isApplicable(message []interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("failed to check message applicability: %v, %s", r, fmt.Sprintf("%s", message))
		}
	}()

	if len(message) < 2 {
		return false
	}

	body := message[2].([]interface{})[0].([]interface{})

	return body[0].(string) == "i"
}

func (p *Parser) Parse(message []interface{}) (*OrderBookDump, error) {
	var dto OrderBookDump

	dto.ChannelId = p.extractChannelId(message)
	dto.SequenceNumber = p.extractSequenceNumber(message)

	payload := p.extractMainPayload(message)

	timestamp, err := p.extractTimestamp(payload)
	if err != nil {
		return nil, err
	}
	dto.Timestamp = timestamp

	dto.CurrencyPair = p.extractCurrencyPair(payload)
	orderBook := p.extractOrderBook(payload)

	asks, err := p.parseBookItems(p.extractAsks(orderBook))
	if err != nil {
		return nil, err
	}
	dto.OrderBook.Asks = asks

	bids, err := p.parseBookItems(p.extractBids(orderBook))
	if err != nil {
		return nil, err
	}
	dto.OrderBook.Bids = bids

	return &dto, nil
}

func (p *Parser) extractChannelId(message []interface{}) int64 {
	return int64(message[0].(float64))
}

func (p *Parser) extractSequenceNumber(message []interface{}) int64 {
	return int64(message[1].(float64))
}

func (p *Parser) extractMainPayload(message []interface{}) []interface{} {
	return message[2].([]interface{})[0].([]interface{})
}

func (p *Parser) extractTimestamp(payload []interface{}) (time.Time, error) {
	timestamp, err := strconv.Atoi(payload[2].(string))
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(timestamp), 0), nil
}

func (p *Parser) extractCurrencyPair(payload []interface{}) string {
	return payload[1].(map[string]interface{})["currencyPair"].(string)
}

func (p *Parser) extractOrderBook(payload []interface{}) []interface{} {
	return payload[1].(map[string]interface{})["orderBook"].([]interface{})
}

func (p *Parser) extractAsks(orderBook []interface{}) map[string]interface{} {
	return orderBook[0].(map[string]interface{})
}

func (p *Parser) extractBids(orderBook []interface{}) map[string]interface{} {
	return orderBook[1].(map[string]interface{})
}

func (p *Parser) parseBookItems(items map[string]interface{}) ([]BookItem, error) {
	result := make([]BookItem, 0, len(items))
	var item BookItem
	for price, size := range items {
		price, err := strconv.ParseFloat(price, 64)
		if err != nil {
			return nil, err
		}
		size, err := strconv.ParseFloat(size.(string), 64)
		if err != nil {
			return nil, err
		}

		item = BookItem{
			Price(price): Size(size),
		}
		result = append(result, item)
	}

	return result, nil
}
