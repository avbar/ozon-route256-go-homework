package sender

import (
	"context"
	"fmt"
	"log"
	"route256/kafka/pkg/order"
	"time"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/encoding/protojson"
)

type Handler func(orderID int64, status string)

type orderSender struct {
	producer  sarama.AsyncProducer
	topic     string
	onSuccess Handler
	onError   Handler
}

func NewOrderSender(producer sarama.AsyncProducer, topic string) *orderSender {
	s := &orderSender{
		producer: producer,
		topic:    topic,
	}

	// config.Producer.Return.Errors = true
	go func() {
		for e := range producer.Errors() {
			key, err := e.Msg.Key.Encode()
			if err != nil {
				log.Printf("kafka message encoding error: %v", err)
				continue
			}
			value, err := e.Msg.Value.Encode()
			if err != nil {
				log.Printf("kafka message encoding error: %v", err)
				continue
			}

			var orderpb order.Order
			err = protojson.Unmarshal(value, &orderpb)
			if err != nil {
				log.Printf("kafka message unmarshalling error: %v", err)
				continue
			}

			if s.onError != nil {
				s.onError(orderpb.GetOrderId(), orderpb.GetStatus())
			}

			log.Printf("order id: %s, status: %s, error: %s", key, orderpb.GetStatus(), e.Error())
		}
	}()

	// config.Producer.Return.Successes = true
	go func() {
		for m := range producer.Successes() {
			key, err := m.Key.Encode()
			if err != nil {
				log.Printf("kafka message encoding error: %v", err)
				continue
			}
			value, err := m.Value.Encode()
			if err != nil {
				log.Printf("kafka message encoding error: %v", err)
				continue
			}

			var orderpb order.Order
			err = protojson.Unmarshal(value, &orderpb)
			if err != nil {
				log.Printf("kafka message unmarshalling error: %v", err)
				continue
			}

			if s.onSuccess != nil {
				s.onSuccess(orderpb.GetOrderId(), orderpb.GetStatus())
			}

			log.Printf("order id: %s, status: %s, partition: %d, offset: %d", key, orderpb.GetStatus(), m.Partition, m.Offset)
		}
	}()

	return s
}

func (s *orderSender) SendOrderStatus(ctx context.Context, orderID int64, status string) {
	orderpb := &order.Order{
		OrderId: orderID,
		Status:  status,
	}

	bytes, err := protojson.Marshal(orderpb)
	if err != nil {
		log.Printf("kafka message marshalling error: %v", err)
		return
	}

	msg := &sarama.ProducerMessage{
		Topic:     s.topic,
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(orderID)),
		Value:     sarama.ByteEncoder(bytes),
		Timestamp: time.Now(),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("order-id"),
				Value: []byte(fmt.Sprint(orderID)),
			},
		},
	}

	s.producer.Input() <- msg
}

func (s *orderSender) AddSuccessHandler(ctx context.Context, onSuccess func(orderID int64, status string)) {
	s.onSuccess = onSuccess
}

func (s *orderSender) AddErrorHandler(ctx context.Context, onError func(orderID int64, status string)) {
	s.onError = onError
}
