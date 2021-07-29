package models

import (
	"testing"
	"time"

	pb "github.com/roman-wb/price-service/internal/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestToPBListReplyPrice(t *testing.T) {
	now := time.Now().UTC()
	price := Price{
		Name:      "Product",
		Price:     100.99,
		Changes:   11,
		UpdatedAt: now,
	}

	want := &pb.ListReply_Price{
		Name:      "Product",
		Price:     100.99,
		Changes:   int64(11),
		UpdatedAt: timestamppb.New(now),
	}

	got := price.ToPBListReplyPrice()

	require.Equal(t, want, got)
}
