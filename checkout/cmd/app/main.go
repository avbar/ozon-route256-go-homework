package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"route256/checkout/internal/api/checkout"
	"route256/checkout/internal/clients/grpc/loms"
	"route256/checkout/internal/clients/grpc/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	repository "route256/checkout/internal/repository/postgres"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/postgres/transactor"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// Clients
	connLOMS, err := grpc.Dial(config.ConfigData.Services.LOMS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to LOMS server: %v", err)
	}
	defer connLOMS.Close()

	connProduct, err := grpc.Dial(config.ConfigData.Services.ProductService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to ProductService server: %v", err)
	}
	defer connLOMS.Close()

	lomsClient := loms.NewClient(connLOMS)
	productClient := productservice.NewClient(connProduct)
	tm := transactor.NewTransactionManager(pool)
	businessLogic := domain.New(lomsClient, productClient, repository.NewCheckoutRepo(tm), tm)

	// Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ConfigData.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterCheckoutServer(s, checkout.NewCheckout(businessLogic))

	log.Printf("checkout server listening at %v port", config.ConfigData.GRPCPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
