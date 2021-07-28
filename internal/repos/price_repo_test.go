package repos_test

import (
	"context"
	"testing"
	"time"

	"github.com/roman-wb/price-service/internal/database"
	"github.com/roman-wb/price-service/internal/models"
	"github.com/roman-wb/price-service/internal/repos"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MongoDB  = "test_price_service"
	MongoURI = "mongodb://localhost:27017/" + MongoDB
)

type PriceRepoTestSuite struct {
	suite.Suite

	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

func (suite *PriceRepoTestSuite) ClearCollection() {
	_, err := suite.collection.DeleteMany(context.Background(), bson.M{}, nil)
	suite.Require().Nil(err)
}

func (suite *PriceRepoTestSuite) SetupTest() {
	client, err := database.NewClient(context.Background(), MongoURI, "file://../../migrations")
	suite.Require().Nil(err)

	suite.client = client
	suite.db = suite.client.Database(MongoDB)
	suite.collection = suite.db.Collection(repos.PriceCollection)

	suite.ClearCollection()
}

func (suite *PriceRepoTestSuite) TearDownSuite() {
	suite.ClearCollection()
}

func (suite *PriceRepoTestSuite) TearDownAllSuite() {
	err := suite.client.Disconnect(context.Background())
	suite.Require().Nil(err)
}

func TestPriceRepo(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, &PriceRepoTestSuite{})
}

func (suite *PriceRepoTestSuite) TestImport() {
	now := time.Now().UTC()
	repo := repos.NewPriceRepo(suite.db)

	testCases := []struct {
		name       string
		prices1    []models.Price
		prices2    []models.Price
		now        time.Time
		wantLen    int
		wantPrices []models.Price
	}{
		{
			name: "Insert to empty collection",
			now:  now,
			prices1: []models.Price{
				{Name: "Product 1", Price: 0},
				{Name: "Product 2", Price: 100.99},
			},
			wantLen: 2,
			wantPrices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now},
			},
		},
		{
			name: "Insert and update collection",
			now:  now,
			prices1: []models.Price{
				{Name: "Product 1", Price: 0},
				{Name: "Product 2", Price: 100.99},
			},
			prices2: []models.Price{
				{Name: "Product 1", Price: 99},
				{Name: "Product 3", Price: 5000},
			},
			wantLen: 3,
			wantPrices: []models.Price{
				{Name: "Product 1", Price: 99, Changes: 2, UpdatedAt: now},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now},
				{Name: "Product 3", Price: 5000, Changes: 1, UpdatedAt: now},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.ClearCollection()

			err := repo.Import(tc.now, tc.prices1)
			suite.Require().Nil(err)

			if len(tc.prices2) > 0 {
				err := repo.Import(tc.now, tc.prices2)
				suite.Require().Nil(err)
			}

			cursor, err := suite.collection.Find(context.Background(), bson.M{}, nil)
			suite.Require().Nil(err)

			var gotPrices []models.Price
			err = cursor.All(context.Background(), &gotPrices)
			suite.Require().Nil(err)

			suite.Require().Equal(tc.wantLen, len(gotPrices))

			for i := range tc.wantPrices {
				suite.Require().NotEmpty(gotPrices[i].ID)
				suite.Require().Equal(tc.wantPrices[i].Name, gotPrices[i].Name)
				suite.Require().Equal(tc.wantPrices[i].Price, gotPrices[i].Price)
				suite.Require().Equal(tc.wantPrices[i].Changes, gotPrices[i].Changes)
				wantDate := tc.wantPrices[i].UpdatedAt.Truncate(time.Second)
				gotDate := gotPrices[i].UpdatedAt.Truncate(time.Second)
				suite.Require().Equal(wantDate, gotDate)
			}
		})
	}
}
