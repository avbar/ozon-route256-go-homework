package main

import (
	"fmt"
	"log"
	"net"
	"route256/loms/internal/api/loms"
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterLOMSServer(s, loms.NewLOMS(domain.New()))

	log.Printf("loms server listening at %v port", grpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
