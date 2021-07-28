package servers

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/roman-wb/price-service/internal/models"
	pb "github.com/roman-wb/price-service/internal/proto"
	"github.com/roman-wb/price-service/internal/servers/mocks"
	"github.com/stretchr/testify/require"
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
		url  string

		mockParserData   []models.Price
		mockParserErr    error
		mockPriceRepoErr error

		wantStatus  string
		wantMessage string
		wantErr     error
	}{
		{
			name:          "Error parser",
			url:           "",
			mockParserErr: errors.New(`parse "": empty url`),
			wantStatus:    "error",
			wantMessage:   `parse "": empty url`,
		},
		{
			name: "Error import",
			url:  "http://yandex.ru",
			mockParserData: []models.Price{
				{Name: "Product 1", Price: 0},
				{Name: "Product 2", Price: 100.99},
			},
			mockPriceRepoErr: errors.New(`some error...`),
			wantStatus:       "error",
			wantMessage:      `some error...`,
		},
		{
			name: "Success import",
			url:  "http://yandex.ru",
			mockParserData: []models.Price{
				{Name: "Product 1", Price: 0},
				{Name: "Product 2", Price: 100.99},
			},
			wantStatus: "ok",
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
				Do(tc.url).
				Return(tc.mockParserData, tc.mockParserErr)

			mockPriceRepo := mocks.NewMockPriceRepo(ctrl)
			if tc.mockPriceRepoErr != nil || tc.mockParserData != nil {
				mockPriceRepo.
					EXPECT().
					Import(gomock.Any(), tc.mockParserData).
					Return(tc.mockPriceRepoErr)
			}

			priceServer := NewPriceServer(mockParser, mockPriceRepo)
			request := &pb.FetchRequest{Url: tc.url}

			gotReply, gotErr := priceServer.Fetch(context.Background(), request)

			require.Equal(t, tc.wantStatus, gotReply.Status)
			require.Equal(t, tc.wantMessage, gotReply.Message)
			require.Nil(t, gotErr)
		})
	}
}
