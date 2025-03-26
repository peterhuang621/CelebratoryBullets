package main

import (
	"context"
	"fmt"
	"log"
	"server/grades"
	"server/mylog"
	"server/registry"
	"server/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	host, port := "localhost", "6000"
	serviceaddr := fmt.Sprintf("http://%v:%v", host, port)

	res := registry.Registration{
		ServiceName:      registry.GradingServie,
		ServiceURL:       serviceaddr,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceaddr + "/services",
		HeartBeatURL:     serviceaddr + "/heartbeat",
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
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at: %s\n", logProvider)
		mylog.SetClientLogger(logProvider, res.ServiceName)
	}

	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
