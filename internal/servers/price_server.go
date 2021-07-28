//go:generate mockgen -destination mocks/price_server.go -package=mocks . Parser,PriceRepo

package servers

import (
	"context"
	"log"
	"time"

	"github.com/roman-wb/price-service/internal/models"
	pb "github.com/roman-wb/price-service/internal/proto"
)

type Parser interface {
	Do(rawurl string) ([]models.Price, error)
}

type PriceRepo interface {
	Import(updatedAt time.Time, prices []models.Price) error
}

type PriceServer struct {
	pb.UnimplementedPriceServer

	parser    Parser
	priceRepo PriceRepo
}

func NewPriceServer(parser Parser, priceRepo PriceRepo) PriceServer {
	return PriceServer{
		parser:    parser,
		priceRepo: priceRepo,
	}
}

func (s *PriceServer) Fetch(ctx context.Context, in *pb.FetchRequest) (*pb.FetchReply, error) {
	log.Printf("Received URL: %v", in.Url)

	// Parse
	prices, err := s.parser.Do(in.Url)
	if err != nil {
		return &pb.FetchReply{Status: "error", Message: err.Error()}, nil
	}

	// Import
	err = s.priceRepo.Import(time.Now().UTC(), prices)
	if err != nil {
		return &pb.FetchReply{Status: "error", Message: err.Error()}, nil
	}

	return &pb.FetchReply{Status: "ok"}, nil
}
