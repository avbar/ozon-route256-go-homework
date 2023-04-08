package main

import (
	"context"
	"log"
	"route256/kafka/kafka"
	"route256/kafka/orders/receiver"
	"route256/notifications/internal/config"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	log.Printf("kafka brokers: %v", config.ConfigData.Brokers)
	consumer, err := kafka.NewConsumer(config.ConfigData.Brokers)
	if err != nil {
		log.Fatalln(err)
	}
	receiver := receiver.NewReceiver(consumer)
	receiver.Subscribe("orders")

	log.Printf("notifications server started")

	<-context.TODO().Done()
}
