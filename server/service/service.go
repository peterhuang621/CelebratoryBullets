package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(ctx context.Context, serviceName, host, port string, registerHanldersFunc func(*gin.Engine), r *gin.Engine) (context.Context, error) {
	registerHanldersFunc(r)
	ctx = startService(ctx, serviceName, host, port, r)
	return ctx, nil
}

func startService(ctx context.Context, serviceName, host, port string, r *gin.Engine) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r.Handler(),
	}
	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()
	go func() {
		fmt.Printf("%v startd, press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
