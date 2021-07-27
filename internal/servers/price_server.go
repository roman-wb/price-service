//go:generate mockgen -destination mock_servers/parser.go . Parser

package servers

import (
	"context"
	"log"

	"github.com/roman-wb/price-service/internal/models"
	pb "github.com/roman-wb/price-service/internal/proto"
)

type Parser interface {
	Do(rawurl string) ([]models.Price, error)
}

type PriceServer struct {
	pb.UnimplementedPriceServer

	parser Parser
}

func NewPriceServer(parser Parser) PriceServer {
	return PriceServer{
		parser: parser,
	}
}

func (s *PriceServer) Fetch(ctx context.Context, in *pb.FetchRequest) (*pb.FetchReply, error) {
	log.Printf("Received URL: %v", in.Url)

	// parse
	data, err := s.parser.Do(in.Url)
	if err != nil {
		return &pb.FetchReply{Status: "error", Message: err.Error()}, nil
	}
	_ = data

	// store
	// repo := NewRepo{}
	// err = repo.Update(data)
	// if err != nil {
	// 	return &pb.FetchReply{Status: "error", Message: err.Error()}, nil
	// }

	return &pb.FetchReply{Status: "ok"}, nil
}
