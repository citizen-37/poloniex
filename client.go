package poloniex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/citizen-37/poloniex/messages"
	"github.com/gorilla/websocket"
)

type (
	Config struct {
		Url            string
		Header         http.Header
		Pair           string
		TimeoutSeconds time.Duration
	}

	Client struct {
		config            Config
		handlers          []messages.Handler
		conn              *websocket.Conn
		closeChan         chan struct{}
		pingChan          chan struct{}
		pingMessage       string
		livelinessChecker *LivelinessChecker
	}
)

func NewClient(config Config, handlers []messages.Handler) *Client {
	return &Client{
		config:      config,
		handlers:    handlers,
		closeChan:   make(chan struct{}),
		pingChan:    make(chan struct{}),
		pingMessage: "poloniex",
	}
}

func (c *Client) Run() error {
	err := c.connect()
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	c.livelinessChecker = newLivelinessChecker(c.conn, c.config.TimeoutSeconds)
	go c.livelinessChecker.watch()

	err = c.subscribe(c.config.Pair)
	if err != nil {
		return fmt.Errorf("subscription failed: %w", err)
	}

	return c.listen()
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) connect() error {
	ws, _, err := websocket.DefaultDialer.Dial(c.config.Url, c.config.Header)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	c.conn = ws

	return nil
}

func (c *Client) subscribe(pair string) error {
	command, err := assembleCommand(pair)
	if err != nil {
		return fmt.Errorf("cannot assemble command: %w", err)
	}

	return c.conn.WriteMessage(websocket.TextMessage, command)
}

func (c *Client) listen() error {
	var message []interface{}

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if e, ok := err.(*websocket.CloseError); ok && e.Code == websocket.CloseNormalClosure {
				return nil
			}
			return fmt.Errorf("failed reading message: %w", err)
		}

		err = json.Unmarshal(data, &message)
		if err != nil {
			return fmt.Errorf("cannot unmarshal message: %w, %s", err, string(data))
		}

		err = c.process(message)
		if err != nil {
			return fmt.Errorf("cannot process message: %w", err)
		}
	}
}

func (c *Client) process(message []interface{}) error {
	for _, handler := range c.handlers {
		if !handler.IsApplicable(message) {
			continue
		}

		err := handler.Process(message)
		if err != nil {
			return fmt.Errorf("failed handling message: %w", err)
		}
	}

	return nil
}

func assembleCommand(pair string) ([]byte, error) {
	return json.Marshal(struct {
		Command string `json:"command"`
		Channel string `json:"channel"`
	}{
		Command: "subscribe",
		Channel: pair,
	})
}
