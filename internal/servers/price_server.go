//go:generate mockgen -destination mocks/price_server.go -package=mocks . Logger,Parser,PriceRepo

package servers

import (
	"context"
	"time"

	"github.com/roman-wb/price-service/internal/models"
	pb "github.com/roman-wb/price-service/internal/proto"
)

type Logger interface {
	Infof(template string, args ...interface{})
}

type Parser interface {
	Fetch(rawurl string) ([]models.Price, error)
}

type PriceRepo interface {
	Import(updatedAt time.Time, prices []models.Price) error
	List(skip int, limit int, orderBy string, orderType int32) ([]models.Price, error)
}

type PriceServer struct {
	pb.UnimplementedPriceServer

	logger    Logger
	parser    Parser
	priceRepo PriceRepo
}

func NewPriceServer(logger Logger, parser Parser, priceRepo PriceRepo) PriceServer {
	return PriceServer{
		logger:    logger,
		parser:    parser,
		priceRepo: priceRepo,
	}
}

func (s *PriceServer) Fetch(ctx context.Context, in *pb.FetchRequest) (*pb.FetchReply, error) {
	s.logger.Infof("Received: %v", in)

	prices, err := s.parser.Fetch(in.Url)
	if err != nil {
		return nil, err
	}

	err = s.priceRepo.Import(time.Now().UTC(), prices)
	if err != nil {
		return nil, err
	}

	return &pb.FetchReply{}, nil
}

func (s *PriceServer) List(ctx context.Context, in *pb.ListRequest) (*pb.ListReply, error) {
	s.logger.Infof("Received: %v", in)

	prices, err := s.priceRepo.List(int(in.Skip), int(in.Limit), in.OrderBy, in.OrderType)
	if err != nil {
		return nil, err
	}

	results := []*pb.ListReply_Price{}
	for _, price := range prices {
		results = append(results, price.ToPBListReplyPrice())
	}

	return &pb.ListReply{Results: results}, nil
}
