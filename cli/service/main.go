package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/roman-wb/price-service/internal/database"
	"github.com/roman-wb/price-service/internal/parser"
	pb "github.com/roman-wb/price-service/internal/proto"
	"github.com/roman-wb/price-service/internal/repos"
	"github.com/roman-wb/price-service/internal/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	Addr     = ":50051"
	MongoDB  = "price_service"
	MongoURI = "mongodb://localhost:27017/" + MongoDB
)

func main() {
	// Mongo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := database.NewClient(ctx, MongoURI, "file://migrations")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	db := client.Database(MongoDB)

	// Deps
	parser := parser.NewParser(&http.Client{})
	priceRepo := repos.NewPriceRepo(db)
	priceServer := servers.NewPriceServer(parser, priceRepo)

	// GRPC
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterPriceServer(grpcServer, &priceServer)

	// Server
	log.Println("Service listen on " + Addr)
	listen, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
