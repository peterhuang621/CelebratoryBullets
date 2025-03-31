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

	server.RegistryServices(engine)
	server.StartKernelServers(&ctx)
	defer server.Close()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", configs.Server_Port),
		Handler: engine,
	}

	go func() {
		defer cancel()
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("server error: %v", err)
			return
		}
	}()

	go func() {
		defer cancel()
		fmt.Print("Enter any key to exit:\n")
		var s string
		fmt.Scanln(&s)
	}()
	<-ctx.Done()
	log.Printf("CTX DONE")
}
