package main

import (
	"context"
	"net"
	"net/http"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"

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
	// Logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Mongo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := database.NewClient(ctx, MongoURI, "file://migrations")
	if err != nil {
		logger.Sugar().Fatalf("failed connection to mongo: %v", err)
	}

	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			logger.Sugar().Fatalf("failed disconnect from mongo: %v", err)
		}
	}()

	db := client.Database(MongoDB)

	// Deps
	parser := parser.NewParser(&http.Client{})
	priceRepo := repos.NewPriceRepo(db)
	priceServer := servers.NewPriceServer(logger.Sugar(), parser, priceRepo)

	// GRPC
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterPriceServer(grpcServer, &priceServer)

	// Server
	logger.Sugar().Infof("Service listen on %s", Addr)
	listen, err := net.Listen("tcp", Addr)
	if err != nil {
		logger.Sugar().Fatalf("failed to listen: %v", err)
	}
	if err := grpcServer.Serve(listen); err != nil {
		logger.Sugar().Fatalf("failed to serve: %v", err)
	}
}
