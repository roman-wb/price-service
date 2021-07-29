package servers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/roman-wb/price-service/internal/models"
	pb "github.com/roman-wb/price-service/internal/proto"
	"github.com/roman-wb/price-service/internal/servers/mocks"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewPriceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wantMockParser := mocks.NewMockParser(ctrl)
	wantMockPriceRepo := mocks.NewMockPriceRepo(ctrl)

	gotPriceServer := NewPriceServer(wantMockParser, wantMockPriceRepo)

	require.NotNil(t, gotPriceServer)
	require.Equal(t, wantMockParser, gotPriceServer.parser)
	require.Equal(t, wantMockPriceRepo, gotPriceServer.priceRepo)
}

func TestPriceServerFetch(t *testing.T) {
	testCases := []struct {
		name string

		url             string
		isMockPriceRepo bool

		mockParserPrices []models.Price
		mockParserErr    error
		mockPriceRepoErr error

		wantReply *pb.FetchReply
		wantErr   error
	}{
		{
			name: "Parser returns error",

			url: "",

			mockParserErr: errors.New(`parse "": empty url`),

			wantReply: nil,
			wantErr:   errors.New(`parse "": empty url`),
		},
		{
			name: "Repo returns error",

			url:             "http://yandex.ru",
			isMockPriceRepo: true,

			mockParserPrices: []models.Price{
				{Name: "Product 1", Price: 0},
				{Name: "Product 2", Price: 100.99},
			},
			mockPriceRepoErr: errors.New(`some error...`),

			wantReply: nil,
			wantErr:   errors.New(`some error...`),
		},
		{
			name: "Response without errors",

			url:             "http://yandex.ru",
			isMockPriceRepo: true,

			mockParserPrices: []models.Price{
				{Name: "Product 1", Price: 0},
				{Name: "Product 2", Price: 100.99},
			},
			mockParserErr:    nil,
			mockPriceRepoErr: nil,

			wantReply: &pb.FetchReply{},
			wantErr:   nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockParser := mocks.NewMockParser(ctrl)
			mockParser.
				EXPECT().
				Fetch(tc.url).
				Return(tc.mockParserPrices, tc.mockParserErr)

			mockPriceRepo := mocks.NewMockPriceRepo(ctrl)
			if tc.isMockPriceRepo {
				mockPriceRepo.
					EXPECT().
					Import(gomock.Any(), tc.mockParserPrices).
					Return(tc.mockPriceRepoErr)
			}

			priceServer := NewPriceServer(mockParser, mockPriceRepo)
			request := &pb.FetchRequest{Url: tc.url}

			gotReply, gotErr := priceServer.Fetch(context.Background(), request)

			require.Equal(t, tc.wantReply, gotReply)
			require.Equal(t, tc.wantErr, gotErr)
		})
	}
}

func TestPriceServerList(t *testing.T) {
	now := time.Now().UTC()

	testCases := []struct {
		name string

		skip      int
		limit     int
		orderBy   string
		orderType int32

		mockPriceRepoPrices []models.Price
		mockPriceRepoErr    error

		wantResults []*pb.ListReply_Price
		wantErr     error
	}{
		{
			name: "Repo returns error",

			skip:      1,
			limit:     100,
			orderBy:   "name",
			orderType: 1,

			mockPriceRepoPrices: []models.Price{},
			mockPriceRepoErr:    errors.New(`some error...`),

			wantResults: nil,
			wantErr:     errors.New(`some error...`),
		},
		{
			name: "Repo returns empty result",

			skip:      1,
			limit:     100,
			orderBy:   "name",
			orderType: 1,

			mockPriceRepoPrices: []models.Price{},
			mockPriceRepoErr:    nil,

			wantResults: nil,
			wantErr:     nil,
		},
		{
			name: "Repo returns results",

			skip:      1,
			limit:     100,
			orderBy:   "name",
			orderType: 1,

			mockPriceRepoPrices: []models.Price{
				{Name: "Product 1", Price: 100.99, Changes: 11, UpdatedAt: now},
				{Name: "Product 2", Price: 0, Changes: 1, UpdatedAt: now},
			},
			mockPriceRepoErr: nil,

			wantResults: []*pb.ListReply_Price{
				{Name: "Product 1", Price: 100.99, Changes: 11, UpdatedAt: timestamppb.New(now)},
				{Name: "Product 2", Price: 0, Changes: 1, UpdatedAt: timestamppb.New(now)},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPriceRepo := mocks.NewMockPriceRepo(ctrl)
			mockPriceRepo.
				EXPECT().
				List(tc.skip, tc.limit, tc.orderBy, tc.orderType).
				Return(tc.mockPriceRepoPrices, tc.mockPriceRepoErr)

			priceServer := NewPriceServer(nil, mockPriceRepo)
			request := &pb.ListRequest{Skip: int64(tc.skip), Limit: int64(tc.limit), OrderBy: tc.orderBy, OrderType: int32(tc.orderType)}

			gotReply, gotErr := priceServer.List(context.Background(), request)

			if len(tc.wantResults) > 0 {
				require.Equal(t, tc.wantResults, gotReply.Results)
			}
			require.Equal(t, tc.wantErr, gotErr)
		})
	}
}
