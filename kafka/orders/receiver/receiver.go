package receiver

import (
	"log"
	"route256/kafka/pkg/order"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/encoding/protojson"
)

type Receiver struct {
	consumer sarama.Consumer
}

func NewReceiver(consumer sarama.Consumer) *Receiver {
	return &Receiver{
		consumer: consumer,
	}
}

func (r *Receiver) Subscribe(topic string) error {
	partitionList, err := r.consumer.Partitions(topic) //get all partitions on the given topic
	if err != nil {
		return err
	}

	for _, partition := range partitionList {
		initialOffset := sarama.OffsetNewest

		pc, err := r.consumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				k := string(message.Key)

				var orderpb order.Order
				err := protojson.Unmarshal(message.Value, &orderpb)
				if err != nil {
					log.Printf("kafka message unmarshalling error: %v", err)
				}

				log.Printf("order id: %s, status: %s, partition: %d, offset: %d",
					k, orderpb.GetStatus(), message.Partition, message.Offset)
			}
		}(pc)
	}

	return nil
}
