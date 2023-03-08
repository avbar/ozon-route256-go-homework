package main

import (
	"fmt"
	"log"
	"net"
	"route256/checkout/internal/api/checkout"
	"route256/checkout/internal/clients/grpc/loms"
	"route256/checkout/internal/clients/grpc/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50050

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	// Clients
	connLOMS, err := grpc.Dial(config.ConfigData.Services.Loms, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	businessLogic := domain.New(lomsClient, productClient)

	// Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterCheckoutServer(s, checkout.NewCheckout(businessLogic))

	log.Printf("checkout server listening at %v port", grpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
