package dump

import (
	"time"
)

type (
	HandlerFunc func(message *OrderBookDump) error

	OrderBookDump struct {
		ChannelId      int64
		SequenceNumber int64
		CurrencyPair   string
		OrderBook      OrderBook
		Timestamp      time.Time
	}

	OrderBook struct {
		Asks []BookItem
		Bids []BookItem
	}

	Price    float64
	Size     float64
	BookItem map[Price]Size
)
