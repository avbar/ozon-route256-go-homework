package receiver

import (
	"route256/kafka/pkg/order"
	"route256/libs/logger"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
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
					logger.Error("kafka message unmarshalling error", zap.Error(err))
				}

				logger.Info("order status changed", zap.String("order id", k), zap.String("status", orderpb.GetStatus()),
					zap.Int32("partition", message.Partition), zap.Int64("offset", message.Offset))
			}
		}(pc)
	}

	return nil
}
