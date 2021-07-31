package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"strings"
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

var addr = flag.String("addr", "localhost:50051", "Listen on host:port")
var mode = flag.String("mode", "dev", "Run mode dev or prod")
var mongo = flag.String("mongo", "mongodb://localhost:27017", "URL to MongoDB without db name")
var dbName = flag.String("dbname", "price_service", "Database name")

func main() {
	flag.Parse()

	// Logger
	var logger *zap.Logger
	if *mode == "prod" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync() //nolint:errcheck

	// Mongo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connURL := strings.TrimRight(*mongo, "/") + "/" + *dbName
	client, err := database.NewClient(ctx, connURL, "file://migrations")
	if err != nil {
		logger.Sugar().Fatalf("failed connection to mongo: %v", err)
	}

	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			logger.Sugar().Fatalf("failed disconnect from mongo: %v", err)
		}
	}()

	db := client.Database(*dbName)

	// Deps
	parser := parser.NewParser(&http.Client{})
	priceRepo := repos.NewPriceRepo(db)
	priceServer := servers.NewPriceServer(logger.Sugar(), parser, priceRepo)

	// GRPC
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterPriceServer(grpcServer, &priceServer)

	// Server
	logger.Sugar().Infof("Service listen on %s", *addr)
	listen, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Sugar().Fatalf("failed to listen: %v", err)
	}
	if err := grpcServer.Serve(listen); err != nil {
		logger.Sugar().Fatalf("failed to serve: %v", err)
	}
}
