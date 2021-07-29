package repos

import (
	"context"
	"time"

	"github.com/roman-wb/price-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const PriceCollection = "prices"

var orderFields = map[string]struct{}{
	"name":       {},
	"price":      {},
	"changes":    {},
	"updated_at": {},
}

type PriceRepo struct {
	collection *mongo.Collection
}

func NewPriceRepo(db *mongo.Database) *PriceRepo {
	return &PriceRepo{
		collection: db.Collection(PriceCollection),
	}
}

func (pr *PriceRepo) Import(updatedAt time.Time, prices []models.Price) error {
	update := []mongo.WriteModel{}
	for _, price := range prices {
		writeModel := pr.updateModel(updatedAt, price)
		update = append(update, writeModel)
	}
	_, err := pr.collection.BulkWrite(context.Background(), update)
	return err
}

func (pr *PriceRepo) List(skip int, limit int, orderBy string, orderType int32) ([]models.Price, error) {
	pipeline := pr.listPipeline(skip, limit, orderBy, orderType)
	cursor, err := pr.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var prices []models.Price
	err = cursor.All(context.Background(), &prices)
	if err != nil {
		return nil, err
	}

	return prices, nil
}

func (pr *PriceRepo) updateModel(updatedAt time.Time, price models.Price) *mongo.UpdateOneModel {
	return mongo.NewUpdateOneModel().
		SetFilter(bson.M{"name": price.Name}).
		SetUpdate(bson.M{
			"$inc": bson.M{
				"changes": 1,
			},
			"$set": bson.M{
				"name":       price.Name,
				"price":      price.Price,
				"updated_at": updatedAt,
			},
		}).
		SetUpsert(true)
}

func (pr *PriceRepo) listPipeline(skip int, limit int, orderBy string, orderType int32) []bson.M {
	if skip < 0 {
		skip = 0
	}

	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	if _, ok := orderFields[orderBy]; !ok {
		orderBy = "name"
	}

	if orderType != -1 {
		orderType = 1
	}

	return []bson.M{
		{"$sort": bson.D{{Key: orderBy, Value: orderType}}},
		{"$skip": skip},
		{"$limit": limit},
	}
}
