package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"server/registry"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Any("/services", (&registry.RegistryService{}).ServeHTTP)
	registry.SetupRegistryService()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := http.Server{
		Addr:    registry.ServerPort,
		Handler: r.Handler(),
	}

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Println("Registry service started. Press any key to stop.")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()
	<-ctx.Done()

	fmt.Println("Shutting down registry service")
}
