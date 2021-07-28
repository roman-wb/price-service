package repos

import (
	"context"
	"time"

	"github.com/roman-wb/price-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const PriceCollection = "prices"

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
		writeModel := pr.importModel(updatedAt, price)
		update = append(update, writeModel)
	}
	_, err := pr.collection.BulkWrite(context.Background(), update)
	return err
}

func (pr *PriceRepo) importModel(updatedAt time.Time, price models.Price) *mongo.UpdateOneModel {
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
