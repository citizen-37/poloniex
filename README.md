Example

```go
package main

import (
	"fmt"
	"github.com/citizen-37/poloniex"
	"github.com/citizen-37/poloniex/messages"
	"github.com/citizen-37/poloniex/messages/dump"
	"github.com/citizen-37/poloniex/messages/update"
	"log"
)

func main() {
	config := poloniex.Config{
		Url:            "wss://api2.poloniex.com",
		Header:         nil,
		Pair:           "USDT_BTC",
		TimeoutSeconds: 10,
	}

	client := poloniex.NewClient(config, []messages.Handler{
		//dump.NewDumpMessageHandler(func(message *dump.OrderBookDump) error {
		//	fmt.Println(message.CurrencyPair)
		//	return nil
		//}),
		update.NewUpdateMessageHandler(func(message *update.OrderBookUpdate) error {
			for _, tx := range message.Transactions {
				fmt.Println(tx)
			}
			return nil
		}),
	})

	done := make(chan struct{})
	go func() {
		err := client.Run()
		if err != nil {
			log.Fatalf("error while reading poloniex: %v", err)
		}

		close(done)
	}()

	<-done
}

```