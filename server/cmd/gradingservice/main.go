package main

import (
	"context"
	"fmt"
	"log"
	"server/grades"
	"server/registry"
	"server/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	host, port := "localhost", "6000"
	serviceaddr := fmt.Sprintf("http://%v:%v", host, port)

	res := registry.Registration{
		ServiceName: registry.GradingServie,
		ServiceURL:  serviceaddr,
	}

	ctx, err := service.Start(context.Background(),
		host,
		port,
		res,
		grades.RegisterHanlers,
		r)
	if err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
