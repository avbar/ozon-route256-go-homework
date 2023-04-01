package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"route256/kafka/kafka"
	"route256/kafka/orders/sender"
	"route256/libs/postgres/transactor"
	"route256/loms/internal/api/loms"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	repository "route256/loms/internal/repository/postgres"
	desc "route256/loms/pkg/loms_v1"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	// DB connection
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, config.ConfigData.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	log.Printf("kafka brokers: %v", config.ConfigData.Brokers)
	producer, err := kafka.NewProducer(config.ConfigData.Brokers)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	sender := sender.NewOrderSender(producer, "orders")

	tm := transactor.NewTransactionManager(pool)
	businessLogic := domain.New(repository.NewLOMSRepo(tm), tm, sender)

	// Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ConfigData.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterLOMSServer(s, loms.NewLOMS(businessLogic))

	log.Printf("loms server listening at %v port", config.ConfigData.GRPCPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
