package main

import (
	"context"
	"fmt"
	"jwt-auth/internal/app"
	serverRouter "jwt-auth/internal/server"
	"log"
	"net/http"
	"time"
)

func main() {
	// parent context;
	ctx := context.Background();

	// start app;
	a, err := app.NewApp(ctx);
	if err != nil {
		log.Fatalf("Application refused to start: %v", err.Error())
		return;
	}

	// shutdown app
	defer func() {
		if err := a.Close(ctx); err != nil {
			log.Printf("Application refused to shutdown: %v", err.Error())
		}
	}()

	// new server router
	r := serverRouter.NewRouter(a);

	// instantiate server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",a.Config.PORT),
		Handler: r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// listen
	log.Printf("server is listening on PORT: %s", a.Config.PORT);

	if err := srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Printf("Server closed!")
			return;
		}
		log.Fatalf("Server Error: %v", err)	
		return
	}
}