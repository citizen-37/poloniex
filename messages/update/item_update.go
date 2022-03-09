package update

import (
	"fmt"
	"strconv"
)

type ItemUpdateParser struct {
}

func NewItemUpdateParser() *ItemUpdateParser {
	return &ItemUpdateParser{}
}

func (p *ItemUpdateParser) IsApplicable(item []interface{}) bool {
	return item[0].(string) == "o"
}

func (p *ItemUpdateParser) ParseAndApply(dto *OrderBookUpdate, item []interface{}) error {

	updateItem, err := p.assembleItem(item)
	if err != nil {
		return fmt.Errorf("failed to assemble update item: %w", err)
	}

	dto.BookUpdates = append(dto.BookUpdates, updateItem)

	return nil
}

func (p *ItemUpdateParser) assembleItem(item []interface{}) (BookUpdate, error) {
	side := p.extractSide(item)

	price, err := p.extractPrice(item)
	if err != nil {
		return BookUpdate{}, err
	}

	size, err := p.extractSize(item)
	if err != nil {
		return BookUpdate{}, err
	}

	epoch, err := p.extractEpochMs(item)
	if err != nil {
		return BookUpdate{}, err
	}

	return BookUpdate{
		Side:    side,
		Price:   price,
		Size:    size,
		EpochMs: epoch,
	}, nil
}

func (p *ItemUpdateParser) extractSide(item []interface{}) string {
	var side string
	if item[1].(float64) > 0 {
		side = bookSideBid
	} else {
		side = bookSideAsk
	}

	return side
}

func (p *ItemUpdateParser) extractPrice(item []interface{}) (float64, error) {
	price, err := strconv.ParseFloat(item[2].(string), 64)
	if err != nil {
		return 0, fmt.Errorf("cannot extract price: %w", err)
	}

	return price, nil
}

func (p *ItemUpdateParser) extractSize(item []interface{}) (float64, error) {
	size, err := strconv.ParseFloat(item[3].(string), 64)
	if err != nil {
		return 0, fmt.Errorf("cannot extract size: %w", err)
	}

	return size, nil
}

func (p *ItemUpdateParser) extractEpochMs(item []interface{}) (int, error) {
	epoch, err := strconv.Atoi(item[4].(string))
	if err != nil {
		return 0, fmt.Errorf("cannot extract epoch_ms: %w", err)
	}

	return epoch, nil
}
