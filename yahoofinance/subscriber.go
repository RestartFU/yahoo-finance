package yahoofinance

import (
	"encoding/base64"
	"errors"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

const yahooFinanceURL = "wss://streamer.finance.yahoo.com/?version=2"

var (
	ErrSubscriberClosed = errors.New("subscriber closed")
	ErrNoDataAvailable  = errors.New("no data available")
)

type Subscriber struct {
	conn    *websocket.Conn
	symbols []Symbol

	mu     sync.Mutex
	queues map[string][]ResponseMessage
	cond   *sync.Cond

	close  sync.Once
	closed bool
}

func NewSubscriber(symbols ...Symbol) *Subscriber {
	sub := &Subscriber{
		symbols: symbols,
		queues:  map[string][]ResponseMessage{},
	}
	sub.cond = sync.NewCond(&sub.mu)
	return sub
}

func (s *Subscriber) Subscribe() error {
	conn, _, err := websocket.DefaultDialer.Dial(yahooFinanceURL, nil)
	if err != nil {
		return err
	}

	var symbols []string
	for _, symbol := range s.symbols {
		symbols = append(symbols, symbol.String())
	}

	err = conn.WriteJSON(SubscribeMessage{
		Subscribe: symbols,
	})
	if err != nil {
		return err
	}
	s.conn = conn

	go s.startPolling()
	return nil
}

func (s *Subscriber) startPolling() {
	defer s.Close()

	for {
		var resp rawResponse
		err := s.conn.ReadJSON(&resp)
		if err != nil || !strings.EqualFold(resp.Type, "pricing") {
			return
		}

		decodedMsg, err := base64.StdEncoding.DecodeString(resp.Message)
		if err != nil {
			return
		}

		var responseMsg ResponseMessage
		err = unmarshalWireMessage(decodedMsg, &responseMsg)
		if err != nil {
			return
		}

		s.mu.Lock()
		s.queues[responseMsg.Symbol] = append(s.queues[responseMsg.Symbol], responseMsg)
		s.mu.Unlock()
		s.cond.Broadcast()
	}
}

func (s *Subscriber) Poll(symbol Symbol) (ResponseMessage, error) {
	if s.closed {
		return ResponseMessage{}, ErrSubscriberClosed
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	symbolString := symbol.String()
	if len(s.queues[symbolString]) == 0 {
		s.cond.Wait()
	}

	response := s.queues[symbolString][0]
	s.queues[symbolString] = s.queues[symbolString][1:]
	return response, nil
}

func (s *Subscriber) Close() error {
	s.close.Do(func() {

	})
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
