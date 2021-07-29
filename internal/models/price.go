package models

import (
	"time"

	pb "github.com/roman-wb/price-service/internal/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Price struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Price     float64            `bson:"price"`
	Changes   int                `bson:"changes"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (p *Price) ToPBListReplyPrice() *pb.ListReply_Price {
	return &pb.ListReply_Price{
		Name:      p.Name,
		Price:     p.Price,
		Changes:   int64(p.Changes),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}
}
