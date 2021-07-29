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
	now1 := time.Now().UTC().Truncate(10 * time.Hour)
	now2 := time.Now().UTC()
	repo := repos.NewPriceRepo(suite.db)

	testCases := []struct {
		name string

		now       time.Time
		oldPrices []models.Price
		newPrices []models.Price

		wantLen    int
		wantPrices []models.Price
	}{
		{
			name: "Insert to empty collection",

			now: now1,
			newPrices: []models.Price{
				{Name: "Product 1", Price: 0},
				{Name: "Product 2", Price: 100.99},
			},

			wantLen: 2,
			wantPrices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now1},
			},
		},
		{
			name: "Insert and update collection",

			now: now2,
			oldPrices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now2},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},
			newPrices: []models.Price{
				{Name: "Product 1", Price: 99},
				{Name: "Product 3", Price: 5000},
			},

			wantLen: 3,
			wantPrices: []models.Price{
				{Name: "Product 1", Price: 99, Changes: 2, UpdatedAt: now2},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
				{Name: "Product 3", Price: 5000, Changes: 1, UpdatedAt: now2},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.ClearCollection()

			for _, price := range tc.oldPrices {
				_, err := suite.collection.InsertOne(context.Background(), price)
				suite.Require().Nil(err)
			}

			err := repo.Import(tc.now, tc.newPrices)
			suite.Require().Nil(err)

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

func (suite *PriceRepoTestSuite) TestList() {
	now1 := time.Now().UTC().Truncate(10 * time.Hour)
	now2 := time.Now().UTC()
	repo := repos.NewPriceRepo(suite.db)

	testCases := []struct {
		name string

		now       time.Time
		skip      int
		limit     int
		orderBy   string
		orderType int32
		prices    []models.Price

		wantPrices []models.Price
		wantErr    error
	}{
		{
			name: "Return empty result",

			wantPrices: nil,
			wantErr:    nil,
		},
		{
			name: "Invalid skip = -1",

			skip: -1,
			prices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},

			wantPrices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},
			wantErr: nil,
		},
		{
			name: "Valid skip = 0",

			skip: 0,
			prices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},

			wantPrices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},
			wantErr: nil,
		},
		{
			name: "Valid skip = 1",

			skip: 1,
			prices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},

			wantPrices: []models.Price{
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},
			wantErr: nil,
		},
		{
			name: "Invalid limit = 0",

			limit: 0,
			prices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},

			wantPrices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},
			wantErr: nil,
		},
		{
			name: "Valid limit = 1",

			limit: 1,
			prices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
				{Name: "Product 2", Price: 100.99, Changes: 1, UpdatedAt: now2},
			},

			wantPrices: []models.Price{
				{Name: "Product 1", Price: 0, Changes: 1, UpdatedAt: now1},
			},
			wantErr: nil,
		},
		{
			name: "Sort by name asc",

			orderBy:   "name",
			orderType: 1,
			prices: []models.Price{
				{Name: "Product 2", Price: 0, Changes: 1, UpdatedAt: now2},
				{Name: "Product 1", Price: 100.99, Changes: 2, UpdatedAt: now1},
			},

			wantPrices: []models.Price{
				{Name: "Product 1", Price: 100.99, Changes: 2, UpdatedAt: now1},
				{Name: "Product 2", Price: 0, Changes: 1, UpdatedAt: now2},
			},
			wantErr: nil,
		},
		{
			name: "Sort by name desc",

			orderBy:   "name",
			orderType: -1,
			prices: []models.Price{
				{Name: "Product 1", Price: 100.99, Changes: 2, UpdatedAt: now1},
				{Name: "Product 2", Price: 0, Changes: 1, UpdatedAt: now2},
			},

			wantPrices: []models.Price{
				{Name: "Product 2", Price: 0, Changes: 1, UpdatedAt: now2},
				{Name: "Product 1", Price: 100.99, Changes: 2, UpdatedAt: now1},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.ClearCollection()

			for _, price := range tc.prices {
				_, err := suite.collection.InsertOne(context.Background(), price)
				suite.Require().Nil(err)
			}

			gotPrices, gotErr := repo.List(tc.skip, tc.limit, tc.orderBy, tc.orderType)

			suite.Require().Equal(len(tc.wantPrices), len(gotPrices))
			for i := range tc.wantPrices {
				suite.Require().Equal(tc.wantPrices[i].Name, gotPrices[i].Name)
				suite.Require().Equal(tc.wantPrices[i].Price, gotPrices[i].Price)
				suite.Require().Equal(tc.wantPrices[i].Changes, gotPrices[i].Changes)
				wantDate := tc.wantPrices[i].UpdatedAt.Truncate(time.Second)
				gotDate := gotPrices[i].UpdatedAt.Truncate(time.Second)
				suite.Require().Equal(wantDate, gotDate)
			}
			suite.Require().Equal(tc.wantErr, gotErr)
		})
	}
}
