package main

import (
	"fmt"
	"log"
	"notes-api/internal/config"
	"notes-api/internal/db"
	"notes-api/internal/server"
)

func main() {
	// load config 
	cfg, err := config.Load();
	if err != nil {
		log.Fatalf("Config error: %v", err);
	}

	// connect to db;
	client, database, err := db.Connect(cfg);
	if(err != nil) {
		log.Fatalf("Connection to DB failed: %v", err);
	}

	// defer db connection;
	defer func() {
		if err := db.Disconnect(client); err != nil {
			log.Printf("Failed to disconnect to DB: %v", err)
		}
	}()

	// new router;
	router := server.NewRouter(database);

	// run 
	addr := fmt.Sprintf(":%s",cfg.ServerPort);

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}