package parser

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/roman-wb/price-service/internal/models"
	"github.com/roman-wb/price-service/internal/parser/mock_parser"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name string
		url  string

		mockBody     string
		mockHttpResp *http.Response
		mockHttpErr  error

		wantData []models.Price
		wantErr  error
	}{
		{
			name:    "Empty URL",
			url:     "",
			wantErr: errors.New(`parse "": empty url`),
		},
		{
			name:    "Invalid URL",
			url:     "yandex.ru/price",
			wantErr: errors.New(`parse "yandex.ru/price": invalid URI for request`),
		},
		{
			name:        "Http error",
			url:         "http://yandex.ru/price",
			mockHttpErr: errors.New("http error..."),
			wantErr:     errors.New(`http error...`),
		},
		{
			name: "Empty data",
			url:  "http://yandex.ru/price",
			mockHttpResp: &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader([]byte(``))),
			},
			wantData: nil,
		},
		{
			name: "Parsed data",
			url:  "http://yandex.ru/price",
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
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var mock *mock_parser.MockHttpClient
			var gotData []models.Price
			var gotErr error

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			if tc.mockHttpResp != nil || tc.mockHttpErr != nil {
				mock = mock_parser.NewMockHttpClient(ctrl)
				mock.
					EXPECT().
					Get(tc.url).
					Return(tc.mockHttpResp, tc.mockHttpErr)
			}

			parser := NewParser(mock)
			gotData, gotErr = parser.Do(tc.url)

			require.Equal(t, len(tc.wantData), len(gotData))
			require.Equal(t, tc.wantData, gotData)
			if tc.wantErr != nil {
				require.Equal(t, tc.wantErr.Error(), gotErr.Error())
			}
		})
	}
}
