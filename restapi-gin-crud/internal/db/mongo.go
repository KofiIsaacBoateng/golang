package db

import (
	"context"
	"fmt"
	"notes-api/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(cfg config.Config) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second);
	defer cancel();

	opts := options.Client().ApplyURI(cfg.MongoURI);

	// connect to db
	client, err := mongo.Connect(opts);
	if err != nil {
		return nil, nil, fmt.Errorf("Mongo DB connection failed: %w", err)
	}

	// checking if client has access to the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, fmt.Errorf("Client PING failed with error: %w", err);
	}

	db := client.Database(cfg.MongoDBName); // create database

	return client, db, nil
}


func Disconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancel();

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("Failed to disconnect: %w", err)
	}

	return nil;
}