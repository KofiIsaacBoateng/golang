package main

import (
	"context"
	"fmt"
	"jwt-auth/internal/app"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	r := gin.Default()

	r.GET("/",  func (c *gin.Context){
			c.JSON(http.StatusOK, gin.H{
				"ok": "true",
				"msg": "Hello jwt auth",
			})})

	// instantiate server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s",a.Config.PORT),
		ReadHeaderTimeout: 5 * time.Second,
		Handler: r,
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