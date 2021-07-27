//go:generate mockgen -destination mock_parser/http_client.go . HttpClient

package parser

import (
	"encoding/csv"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/roman-wb/price-service/internal/models"
)

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Parser struct {
	httpClient HttpClient
}

func NewParser(httpClient HttpClient) *Parser {
	return &Parser{
		httpClient: httpClient,
	}
}

func (p *Parser) Do(rawurl string) ([]models.Price, error) {
	// validate
	_, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}

	// request
	resp, err := p.httpClient.Get(rawurl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return p.Parse(resp.Body)
}

func (p *Parser) Parse(body io.ReadCloser) ([]models.Price, error) {
	var prices []models.Price

	reader := csv.NewReader(body)
	reader.Comma = ';'
	reader.FieldsPerRecord = 2

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		name := strings.TrimSpace(record[0])
		price, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}

		prices = append(prices, models.Price{
			Name:  name,
			Price: price,
		})
	}

	return prices, nil
}