package main

import (
	"flag"
	"fmt"
	"net/http"
	"route256/kafka/kafka"
	"route256/kafka/orders/receiver"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/tracing"
	"route256/notifications/internal/config"

	"go.uber.org/zap"
)

var develMode = flag.Bool("devel", true, "development mode")

func main() {
	flag.Parse()

	logger.Init(*develMode)
	log := logger.GlobalLogger()
	tracing.Init(log, "notifications")

	err := config.Init()
	if err != nil {
		log.Fatal("config init", zap.Error(err))
	}

	log.Info("connect to kafka", zap.Strings("brokers", config.ConfigData.Brokers))
	consumer, err := kafka.NewConsumer(config.ConfigData.Brokers)
	if err != nil {
		log.Fatal("failed to create consumer", zap.Error(err))
	}
	receiver := receiver.NewReceiver(consumer)
	receiver.Subscribe("orders")

	log.Info("notifications server started")

	http.Handle("/metrics", metrics.New())
	http.ListenAndServe(fmt.Sprintf(":%d", config.ConfigData.HTTPPort), nil)
}
