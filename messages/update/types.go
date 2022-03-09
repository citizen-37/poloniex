package update

import (
	"time"
)

const transactionSideBuy string = "buy"
const transactionSideSell string = "sell"
const bookSideAsk string = "ask"
const bookSideBid string = "bid"

type (
	HandlerFunc func(message *OrderBookUpdate) error

	ItemParser interface {
		ParseAndApply(dto *OrderBookUpdate, item []interface{}) error
		IsApplicable(item []interface{}) bool
	}

	OrderBookUpdate struct {
		ChannelId      int64
		SequenceNumber int64
		BookUpdates    []BookUpdate
		Transactions   []Transaction
	}

	BookUpdate struct {
		Side    string
		Price   float64
		Size    float64
		EpochMs int
	}

	Transaction struct {
		Id        int
		Price     float64
		Size      float64
		Side      string
		EpochMs   int
		Timestamp time.Time
	}
)
