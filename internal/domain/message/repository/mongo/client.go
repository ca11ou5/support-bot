package mongo

import (
	"context"
	"github.com/ca11ou5/support-bot/internal/domain/message/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"log/slog"
	"os"
	"time"
)

type Client struct {
	statsColl *mongo.Collection
}

func NewClient(connString string) *Client {
	return &Client{
		statsColl: connectToMongo(connString),
	}
}

func connectToMongo(connString string) *mongo.Collection {
	clientOpts := options.Client().ApplyURI(connString)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	return client.Database("telegram_bot").Collection("stats")
}

func (c *Client) SaveStats(stats entity.Stats) error {
	_, err := c.statsColl.InsertOne(context.Background(), stats)
	return err
}

func (c *Client) GetStats() []entity.Stats {
	opts := options.Find()
	opts.SetLimit(10)
	opts.SetSort(bson.M{"timestamp": -1})

	cur, err := c.statsColl.Find(context.Background(), bson.D{}, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	var results []entity.Stats
	if err = cur.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}

	return results
}
