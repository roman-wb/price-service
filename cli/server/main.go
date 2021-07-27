package main

import (
	"log"
	"net"
	"net/http"

	"github.com/roman-wb/price-service/internal/parser"
	pb "github.com/roman-wb/price-service/internal/proto"
	"github.com/roman-wb/price-service/internal/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	addr = ":50051"
)

func main() {
	// Deps
	parser := parser.NewParser(&http.Client{})
	priceServer := servers.NewPriceServer(parser)

	// GRPC
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterPriceServer(grpcServer, &priceServer)

	// Server
	log.Println("Listen on " + addr)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
