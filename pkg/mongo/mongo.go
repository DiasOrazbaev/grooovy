package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"grovo/config"
	"time"
)

func NewMongoDatabase(cfg *config.Config) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.Mongo.Url))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.Mongo.Database), nil
}
