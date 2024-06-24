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
	wordsColl *mongo.Collection
}

func NewClient(connString string) *Client {
	stats, words := connectToMongo(connString)

	return &Client{
		statsColl: stats,
		wordsColl: words,
	}
}

func connectToMongo(connString string) (*mongo.Collection, *mongo.Collection) {
	clientOpts := options.Client().ApplyURI(connString)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	wordsColl := client.Database("telegram_bot").Collection("words")
	var result map[string]interface{}
	err = wordsColl.FindOne(context.Background(), bson.D{}).Decode(&result)
	if err != nil {
		// Если документ не найден, вставляем новый
		if err == mongo.ErrNoDocuments {
			_, err = wordsColl.InsertOne(context.Background(), bson.M{
				"покупки": 1,
			})
		}
	}

	return client.Database("telegram_bot").Collection("stats"), wordsColl
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

func (c *Client) InsertWords(words map[string]int) error {
	var result map[string]interface{}
	err := c.wordsColl.FindOne(context.Background(), bson.D{}).Decode(&result)
	if err != nil {
		// Если документ не найден, вставляем новый
		if err == mongo.ErrNoDocuments {
			_, err := c.wordsColl.InsertOne(context.Background(), words)
			return err
		}
		return err
	}

	updatedDoc := result
	for k, v := range words {
		if val, ok := updatedDoc[k]; ok {
			if intVal, ok := val.(int32); ok {
				updatedDoc[k] = int(intVal) + v
			} else if intVal, ok := val.(int64); ok {
				updatedDoc[k] = int(intVal) + v
			} else if floatVal, ok := val.(float64); ok {
				updatedDoc[k] = int(floatVal) + v
			} else {
				updatedDoc[k] = v
			}
		} else {
			updatedDoc[k] = v
		}
	}

	filter := bson.M{"_id": updatedDoc["_id"]}
	_, err = c.wordsColl.ReplaceOne(context.Background(), filter, updatedDoc)
	if err != nil {
		return err
	}

	return nil
}
