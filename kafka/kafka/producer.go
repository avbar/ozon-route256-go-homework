package kafka

import (
	"github.com/Shopify/sarama"
)

func NewProducer(brokers []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewHashPartitioner // sarama.NewRoundRobinPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true // exactly once
	config.Net.MaxOpenRequests = 1    // Idempotent producer requires Net.MaxOpenRequests to be 1

	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
