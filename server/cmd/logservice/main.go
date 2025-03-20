package main

import (
	"context"
	"fmt"
	"log"
	"server/mylog"
	"server/service"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func main() {
	r = gin.Default()

	mylog.Run("./distributed.log")
	host, port := "localhost", "4000"
	ctx, err := service.Start(context.Background(), "Log Service", host, port,
		mylog.RegisterHandlers, r)
	if err != nil {
		log.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("shutting down log service.")
}
