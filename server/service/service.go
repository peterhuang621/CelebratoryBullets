package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"server/registry"

	"github.com/gin-gonic/gin"
)

func Start(ctx context.Context, host, port string, reg registry.Registration, registerHanldersFunc func(*gin.Engine), r *gin.Engine) (context.Context, error) {
	registerHanldersFunc(r)
	ctx = startService(ctx, reg.ServiceName, host, port, r)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string, r *gin.Engine) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r.Handler(),
	}
	go func() {
		log.Println(srv.ListenAndServe())
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()
	go func() {
		fmt.Printf("%v startd, press any key to stop. \n", serviceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		cancel()
	}()
	return ctx
}
