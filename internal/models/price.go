package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Price struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Price     float64            `bson:"price"`
	Changes   int                `bson:"changes"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
