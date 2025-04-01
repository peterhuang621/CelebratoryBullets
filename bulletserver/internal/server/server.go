package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/configs"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/pkg"
	goredis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var file *os.File
var err error
var fileMutex *sync.Mutex
var grpc_server *grpc.Server
var redisClient *goredis.Client

func init() {
	fileInfo, err := os.Stat(configs.GL_Drawingfile)
	if os.IsNotExist(err) {
		_, err = os.Create(configs.GL_Drawingfile)
		if err != nil {
			log.Fatalf("Failed on creating the nonexisted file: %v", err)
			return
		}
	}
	file, err = os.OpenFile(configs.GL_Drawingfile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Failed on opening the existed file: %v", err)
		return
	}
	if fileInfo.Size() != 0 {
		startDrawing()
	}
	fileMutex = &sync.Mutex{}
	log.Printf("File successfully opened!\n")
	initMQ()
	log.Printf("initMQ!\n")
}

func StartKernelServers(ctx *context.Context) {
	initgRPC()
	log.Println("init gRPC!")
	initRedis()
	log.Println("init Redis!")
	go StartMQ(ctx)
	log.Println("All kernel servers started!")
}

func RegistryServices(engine *gin.Engine) {
	engine.POST(configs.Client_Addr, func(ctx *gin.Context) {
		var bullets []configs.Bullet
		if err := ctx.ShouldBindJSON(&bullets); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invaild bullets data received"})
			pkg.WarnHandle(err, "invaild bullets data received")
			return
		}
		if err := ValidateBullets(&bullets); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invaild bullets contents: %s", err.Error())})
			pkg.WarnHandle(err)
			return
		}
		err := isAllowed(len(bullets), configs.Redis_Key_DefaultExpr)
		if err != nil {
			ctx.JSON(http.StatusInsufficientStorage, gin.H{"inner_error": err.Error()})
			return
		}

		sendingdata, err := json.Marshal(&bullets)
		if err != nil {
			pkg.WarnHandle(err, "Failed on Marshaling bullets into json format")
			ctx.JSON(http.StatusInternalServerError, gin.H{"inner_error": err.Error()})
			return
		}

		resp, err := http.Post(fmt.Sprintf("http://localhost:%s/%s", configs.Server_Port, configs.Server_Addr), "application/json", bytes.NewBuffer(sendingdata))
		if err != nil {
			pkg.WarnHandle(err, "Failed on Posting to bullets(mq) server")
			ctx.JSON(http.StatusInternalServerError, gin.H{"inner_error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			ans := fmt.Errorf("bullets(mq) server returned status: %d, body: %s", resp.StatusCode, string(body))
			pkg.WarnHandle(ans)
			ctx.JSON(http.StatusInternalServerError, gin.H{"inner_err": ans.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":   "bullets successfully received!",
			"count": len(bullets),
		})
	})

	engine.GET(configs.Client_Addr, func(ctx *gin.Context) {
		currCount, err := getKey()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"inner_error": err.Error()})
		} else {
			ctx.JSON(http.StatusOK, gin.H{configs.Redis_Key: currCount})
		}
	})

	engine.POST(fmt.Sprintf("/%s", configs.Server_Addr), func(ctx *gin.Context) {
		var bullets []configs.Bullet
		if err := ctx.ShouldBindJSON(&bullets); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invaild bullets data received"})
			pkg.WarnHandle(err, "invaild bullets data received")
			return
		}
		SendToMQwithoutRoutingKey(&bullets)
		ctx.JSON(http.StatusOK, nil)
	})

	engine.GET(fmt.Sprintf("/%s", configs.Server_Addr), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"mqaddress": configs.MQ_Addr,
			"rabbitmq":  mqconn != nil,
			"queues":    configs.MQ_QueueNumber,
		})
	})

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, nil)
	})
}

func ValidateBullets(bullets *[]configs.Bullet) error {
	for i, bullet := range *bullets {
		if bullet.DurationSecs <= 0.0 || bullet.DurationSecs > configs.DurationSecs_Max {
			return fmt.Errorf("invalid bullet[%d].DurationSecs (%f)", i, bullet.DurationSecs)
		}

		if bullet.Size <= 0.0 || bullet.Size > configs.Size_Max {
			return fmt.Errorf("invalid bullet[%d].Size (%f)", i, bullet.Size)
		}

		for j, v := range bullet.Color {
			if v < 0 || v > 1 {
				return fmt.Errorf("bullet[%d].Color[%d] (%f)", i, j, v)
			}
		}

		for j, v := range bullet.Position {
			if v > configs.Position_Max {
				return fmt.Errorf("bullet[%d].Position[%d] (%f)", i, j, v)
			}
		}
	}
	return nil
}

func WriteToDrawingFile(bullets *[]configs.Bullet, queuename string) {
	if len(*bullets) == 0 {
		log.Printf("Empty bullets, no writing!\n")
		return
	}

	err := consumeKey(len(*bullets))
	if err != nil {
		pkg.WarnHandle(err, "Failed to consume bullets in Redis server, please check and abort these bullets...")
		return
	}
	log.Printf("Successfully DecrBy in Redis server\n")

	fileMutex.Lock()
	for _, v := range *bullets {
		fmt.Fprintf(file, "%s %s\n", v.String(), queuename)
	}
	fileMutex.Unlock()

	log.Printf("Wrote %d bullet(s) to the file!\n", len(*bullets))
}

func startDrawing() {
	log.Println("Drawing to OpenGL...")
}

func Close() {
	log.Println("Close all things...")
	file.Sync()
	if err = file.Close(); err != nil {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			err = file.Close()
			if err != nil {
				pkg.WarnHandle(err, fmt.Sprintf("Closing the file not properly, try %d times...", i))
			} else {
				break
			}
		}
	}
	pkg.WarnHandle(err, "Failed on closing the file, please manually check the drawing file status")

	grpc_server.GracefulStop()
	log.Println("Grpc_server gracefully ended!")

	if redisClient != nil {
		err := redisClient.Close()
		pkg.WarnHandle(err, "Error while closing Redis client")
		log.Println("Redis server ended!")
	} else {
		log.Println("No Redis server needs to end!")
	}
	log.Println("Close function ended! Program ended!")
}
