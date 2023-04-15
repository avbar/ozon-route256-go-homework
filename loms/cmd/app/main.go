package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"route256/kafka/kafka"
	"route256/kafka/orders/sender"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/postgres/dbwrapper"
	"route256/libs/postgres/transactor"
	"route256/libs/tracing"
	"route256/loms/internal/api/loms"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	repository "route256/loms/internal/repository/postgres"
	desc "route256/loms/pkg/loms_v1"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var develMode = flag.Bool("devel", true, "development mode")

func main() {
	flag.Parse()

	logger.Init(*develMode)
	log := logger.GlobalLogger()
	tracing.Init(log, "loms")

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

	log.Info("connect to kafka", zap.Strings("brokers", config.ConfigData.Brokers))
	producer, err := kafka.NewProducer(config.ConfigData.Brokers)
	if err != nil {
		log.Fatal("failed to create producer", zap.Error(err))
	}
	sender := sender.NewOrderSender(producer, "orders")

	tm := transactor.NewTransactionManager(dbwrapper.NewWrapper(pool))
	businessLogic := domain.New(repository.NewLOMSRepo(tm), tm, sender)

	// Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ConfigData.GRPCPort))
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				logger.LoggingInterceptor,
				otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
				metrics.MetricsInterceptor,
			),
		),
	)
	reflection.Register(s)

	desc.RegisterLOMSServer(s, loms.NewLOMS(businessLogic))

	log.Info("loms server listening", zap.Int("port", config.ConfigData.GRPCPort))

	go func() {
		http.Handle("/metrics", metrics.New())
		http.ListenAndServe(fmt.Sprintf(":%d", config.ConfigData.HTTPPort), nil)
	}()

	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
