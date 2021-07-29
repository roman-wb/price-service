package parser

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/roman-wb/price-service/internal/models"
	"github.com/roman-wb/price-service/internal/parser/mocks"
	"github.com/stretchr/testify/require"
)

func TestParserFetch(t *testing.T) {
	testCases := []struct {
		name string

		url string

		mockHttpResp *http.Response
		mockHttpErr  error

		wantData []models.Price
		wantErr  error
	}{
		{
			name: "Empty URL",

			url: "",

			wantData: nil,
			wantErr:  errors.New(`parse "": empty url`),
		},
		{
			name: "Invalid URL",

			url: "yandex.ru/price",

			wantData: nil,
			wantErr:  errors.New(`parse "yandex.ru/price": invalid URI for request`),
		},
		{
			name: "Http request returns error",

			url: "http://yandex.ru/price",

			mockHttpErr: errors.New("http error..."),

			wantData: nil,
			wantErr:  errors.New(`http error...`),
		},
		{
			name: "Http request returns empty data",

			url: "http://yandex.ru/price",

			mockHttpResp: &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader([]byte(``))),
			},

			wantData: nil,
			wantErr:  nil,
		},
		{
			name: "Parsed data",

			url: "http://yandex.ru/price",

			mockHttpResp: &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`
					Product 1;-1
					Product 2;0
					Product 3;0.99
					Product 4;error
					Product 4;100.99
				`))),
			},

			wantData: []models.Price{
				{Name: "Product 1", Price: -1},
				{Name: "Product 2", Price: 0},
				{Name: "Product 3", Price: 0.99},
				{Name: "Product 4", Price: 100.99},
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

			var mockHttpClient *mocks.MockHttpClient
			if tc.mockHttpResp != nil || tc.mockHttpErr != nil {
				mockHttpClient = mocks.NewMockHttpClient(ctrl)
				mockHttpClient.
					EXPECT().
					Get(tc.url).
					Return(tc.mockHttpResp, tc.mockHttpErr)
			}

			parser := NewParser(mockHttpClient)

			gotData, gotErr := parser.Fetch(tc.url)

			require.Equal(t, tc.wantData, gotData)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr.Error(), gotErr.Error())
			}
		})
	}
}
