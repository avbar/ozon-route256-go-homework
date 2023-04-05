package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"route256/checkout/internal/api/checkout"
	"route256/checkout/internal/clients/grpc/loms"
	"route256/checkout/internal/clients/grpc/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	repository "route256/checkout/internal/repository/postgres"
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/logger"
	"route256/libs/postgres/transactor"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var develMode = flag.Bool("devel", true, "development mode")

func main() {
	flag.Parse()

	log := logger.Init(*develMode)

	err := config.Init()
	if err != nil {
		log.Fatal("config init", zap.Error(err))
	}

	// DB connection
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, config.ConfigData.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to DB", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to ping DB", zap.Error(err))
	}

	// Clients
	connLOMS, err := grpc.Dial(config.ConfigData.Services.LOMS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to LOMS server", zap.Error(err))
	}
	defer connLOMS.Close()

	connProduct, err := grpc.Dial(config.ConfigData.Services.ProductService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to ProductService server", zap.Error(err))
	}
	defer connLOMS.Close()

	lomsClient := loms.NewClient(connLOMS)
	productClient := productservice.NewClient(connProduct)
	tm := transactor.NewTransactionManager(pool)
	businessLogic := domain.New(lomsClient, productClient, repository.NewCheckoutRepo(tm), tm)

	// Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ConfigData.GRPCPort))
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterCheckoutServer(s, checkout.NewCheckout(businessLogic))

	log.Info("checkout server listening", zap.String("port", fmt.Sprint(config.ConfigData.GRPCPort)))

	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
