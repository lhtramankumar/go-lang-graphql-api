// database/db.go

package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client
var db *mongo.Database

func ConnectDB(url string, dbName string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	dbClient = client
	db = client.Database(dbName)

	log.Println("Connected to MongoDB!")
	return nil
}

// GetDBClient returns the MongoDB client instance.
func GetDBClient() *mongo.Client {
	return dbClient
}

// GetDB returns the MongoDB database instance.
func GetDB() *mongo.Database {
	return db
}

type Book struct {
	ID          primitive.ObjectID `json:"_id"`
	Title       string             `json:"title"`
	Bookname    string             `json:"bookname"`
	Description string             `json:"description"`
	Author      string             `json:"author"`
	AddedOn     float64            `json:"addedOn"`
	// Other fields...
}
