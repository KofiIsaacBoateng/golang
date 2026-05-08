package db

import (
	"context"
	"fmt"
	"jwt-auth/internal/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	Database *mongo.Database
}


func Connect(ctx context.Context, cfg config.Config) (*MongoDB, error) {
	// set uri to db options
	opts := options.Client().ApplyURI(cfg.MONGO_URI);
	
	// connect to database
	client, err := mongo.Connect(opts);
	if err != nil {
		return nil, fmt.Errorf("Connection to DB failed: %w", err)
	}

	// ping to verify if deployment is reachable
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Client PING unable to REACH deployment: %w", err)
	}

	// create database cluster
	database := client.Database(cfg.MONGO_DB_NAME);

	return &MongoDB{Client: client, Database: database}, nil
}

func Disconnect(ctx context.Context, client *mongo.Client) error {
	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("Database disconnection FAILED: %w", err)
	}

	return nil;
}