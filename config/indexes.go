package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes creates necessary database indexes
func CreateIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create unique index on invoice_id for invoices collection
	invoicesCollection := db.Collection("invoices")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "invoice_id", Value: 1}, // 1 for ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := invoicesCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Warning: Failed to create invoice_id index: %v", err)
		return err
	}

	// Create unique index on user_id for users collection
	usersCollection := db.Collection("users")

	userIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err = usersCollection.Indexes().CreateOne(ctx, userIndexModel)
	if err != nil {
		log.Printf("Warning: Failed to create user_id index: %v", err)
		return err
	}

	// Create unique index on email for users collection
	emailIndexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	_, err = usersCollection.Indexes().CreateOne(ctx, emailIndexModel)
	if err != nil {
		log.Printf("Warning: Failed to create email index: %v", err)
		return err
	}

	log.Println("Database indexes created successfully!")
	return nil
}
