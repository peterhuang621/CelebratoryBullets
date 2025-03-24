package main

import (
	"context"
	"fmt"
	"log"
	"server/mylog"
	"server/registry"
	"server/service"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func main() {
	r = gin.Default()

	mylog.Run("./distributed.log")
	host, port := "localhost", "4000"
	serviceAddr := fmt.Sprintf("http://%s:%s", host, port)

	res := registry.Registration{
		ServiceName:      registry.LogService,
		ServiceURL:       serviceAddr,
		RequiredServices: make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddr + "/services",
	}

	ctx, err := service.Start(context.Background(), host, port, res, mylog.RegisterHandlers, r)
	if err != nil {
		log.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("shutting down log service.")
}
