package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/configs"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/internal/server"
)

var engine *gin.Engine

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	engine = gin.Default()
	// server.StartMQ(&ctx, &cancel)
	server.SeverServices(engine)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", configs.Server_Port),
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server error: %v", err)
		}
		cancel()
	}()

	go func() {
		fmt.Print("Enter any key to exit:\n")
		var s string
		fmt.Scanln(&s)
		cancel()
	}()
	<-ctx.Done()
	defer server.CloseDrawingFile()
}
