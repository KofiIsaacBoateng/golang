package app

import (
	"context"
	"fmt"
	"jwt-auth/internal/config"
	"jwt-auth/internal/db"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
);

type App struct {
	Config config.Config

	MongoClient *mongo.Client

	MongoDatabase *mongo.Database
}

func NewApp(ctx context.Context) (*App, error) {
	// load configs
	cfg, err := config.Load();
	if err != nil {
		return nil, err;
	} 

	// connect to database;
	db, err := db.Connect(ctx, cfg);
	if err != nil {
		return nil, err;
	}


	return &App{
		Config: cfg,
		MongoClient: db.Client,
		MongoDatabase: db.Database,
	}, nil

}


func (a *App) Close(ctx context.Context) error {
	// if client is missing, connection is already closed!
	if a.MongoClient == nil {
		return nil
	}

	// app shutdown context
	closeCtx, cancel := context.WithTimeout(ctx, 5*time.Second);
	defer cancel()

	// disconnect db
	if err :=db.Disconnect(closeCtx, a.MongoClient); err != nil {
		return fmt.Errorf("App shutdown failed: %w", err);
	}

	return nil
}