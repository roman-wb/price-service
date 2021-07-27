package servers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/roman-wb/price-service/internal/models"
	"github.com/roman-wb/price-service/internal/parser"
	pb "github.com/roman-wb/price-service/internal/proto"
	"github.com/roman-wb/price-service/internal/servers/mock_servers"
	"github.com/stretchr/testify/assert"
)

func TestNewPriceServer(t *testing.T) {
	httpClient := &http.Client{}
	parser := parser.NewParser(httpClient)
	got := NewPriceServer(parser)

	assert.NotNil(t, got)
	assert.Equal(t, got.parser, parser)
}

func TestPriceServerFetch(t *testing.T) {
	testCases := []struct {
		name string
		url  string

		mockParserData []models.Price
		mockParserErr  error

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
			name:        "Store data",
			url:         "http://yandex.ru",
			wantStatus:  "ok",
			wantMessage: ``,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mock_servers.NewMockParser(ctrl)
			mock.
				EXPECT().
				Do(tc.url).
				Return(tc.mockParserData, tc.mockParserErr)

			priceServer := NewPriceServer(mock)
			request := &pb.FetchRequest{Url: tc.url}
			gotReply, gotErr := priceServer.Fetch(context.Background(), request)

			assert.Equal(t, tc.wantStatus, gotReply.Status)
			assert.Equal(t, tc.wantMessage, gotReply.Message)
			assert.Nil(t, gotErr)
		})
	}
}
