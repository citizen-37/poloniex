package update

import (
	"fmt"
	"strconv"
	"time"
)

type ItemTransactionParser struct {
}

func NewItemTransactionParser() *ItemTransactionParser {
	return &ItemTransactionParser{}
}

func (p *ItemTransactionParser) IsApplicable(item []interface{}) bool {
	return item[0].(string) == "t"
}

func (p *ItemTransactionParser) ParseAndApply(dto *OrderBookUpdate, item []interface{}) error {
	tx, err := p.assembleItem(item)
	if err != nil {
		return fmt.Errorf("cannot assemble transaction item: %w", err)
	}

	dto.Transactions = append(dto.Transactions, tx)

	return nil
}

func (p *ItemTransactionParser) assembleItem(item []interface{}) (Transaction, error) {
	side := p.extractSide(item)

	id, err := p.extractId(item)
	if err != nil {
		return Transaction{}, err
	}

	price, err := p.extractPrice(item)
	if err != nil {
		return Transaction{}, err
	}

	size, err := p.extractSize(item)
	if err != nil {
		return Transaction{}, err
	}

	epoch, err := p.extractEpochMs(item)
	if err != nil {
		return Transaction{}, err
	}

	timestamp := p.extractTimestamp(item)

	return Transaction{
		Id:        id,
		Price:     price,
		Size:      size,
		Side:      side,
		Timestamp: timestamp,
		EpochMs:   epoch,
	}, nil
}

func (p *ItemTransactionParser) extractId(item []interface{}) (int, error) {
	id, err := strconv.Atoi(item[1].(string))
	if err != nil {
		return 0, fmt.Errorf("cannot extract id: %w", err)
	}

	return id, nil
}

func (p *ItemTransactionParser) extractSide(item []interface{}) string {
	var side string
	if item[2].(float64) > 0 {
		side = transactionSideBuy
	} else {
		side = transactionSideSell
	}

	return side
}

func (p *ItemTransactionParser) extractPrice(item []interface{}) (float64, error) {
	price, err := strconv.ParseFloat(item[3].(string), 64)
	if err != nil {
		return 0, fmt.Errorf("cannot extract price: %w", err)
	}

	return price, nil
}

func (p *ItemTransactionParser) extractSize(item []interface{}) (float64, error) {
	size, err := strconv.ParseFloat(item[4].(string), 64)
	if err != nil {
		return 0, fmt.Errorf("cannot extract size: %w", err)
	}

	return size, nil
}

func (p *ItemTransactionParser) extractTimestamp(item []interface{}) time.Time {
	return time.Unix(int64(item[5].(float64)), 0)
}

func (p *ItemTransactionParser) extractEpochMs(item []interface{}) (int, error) {
	epoch, err := strconv.Atoi(item[6].(string))
	if err != nil {
		return 0, fmt.Errorf("cannot extract epoch_ms: %w", err)
	}

	return epoch, nil
}
