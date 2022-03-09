package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type LivelinessChecker struct {
	pingChan       chan struct{}
	pingMessage    string
	conn           *websocket.Conn
	timeoutSeconds time.Duration
	ticker         *time.Ticker
}

func newLivelinessChecker(conn *websocket.Conn, timeoutSeconds time.Duration) *LivelinessChecker {
	return &LivelinessChecker{
		pingMessage:    "poloniex",
		pingChan:       make(chan struct{}),
		conn:           conn,
		timeoutSeconds: timeoutSeconds,
		ticker:         time.NewTicker(3 * time.Second),
	}
}

func (l *LivelinessChecker) watch() {
	l.setControlMessageHandlers()

	lastPingAt := time.Now()
	for {
		select {
		case <-l.ticker.C:
			err := l.ping()
			if err != nil {
				log.Fatalf("cannot send ping message: %v", err)
				return
			}
		case <-l.pingChan:
			lastPingAt = time.Now()
		default:
			if time.Now().Sub(lastPingAt) > l.timeoutSeconds*time.Second {
				l.die()
				return
			}
		}
	}
}

func (l *LivelinessChecker) die() {
	l.ticker.Stop()

	err := l.conn.Close()
	if err != nil {
		panic(fmt.Sprintf("error closing timedout connection: %v", err))
	}
}

func (l *LivelinessChecker) setControlMessageHandlers() {
	l.conn.SetCloseHandler(func(code int, text string) error {
		l.die()
		return nil
	})

	l.conn.SetPingHandler(func(appData string) error {
		return l.conn.WriteMessage(websocket.PongMessage, []byte(appData))
	})

	l.conn.SetPongHandler(func(appData string) error {
		if appData == l.pingMessage {
			l.pingChan <- struct{}{}
		}
		return nil
	})
}

func (l *LivelinessChecker) ping() error {
	return l.conn.WriteMessage(websocket.PingMessage, []byte(l.pingMessage))
}
