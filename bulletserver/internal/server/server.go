package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/configs"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/pkg"
)

var file *os.File
var err error
var fileMutex *sync.Mutex

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
	StartMQ()
}

func SeverServices(engine *gin.Engine) {
	engine.POST(fmt.Sprintf("/%s", configs.Server_Addr), func(ctx *gin.Context) {
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
		SendToMQ(&bullets)
		ctx.JSON(http.StatusOK, gin.H{
			"msg":   "bullets successfully received!",
			"count": len(bullets),
		})
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
		if bullet.DurationSecs <= 0 || bullet.DurationSecs > configs.DurationSecs_Max {
			return fmt.Errorf("invalid bullet[%d].DurationSecs (%d)", i, bullet.DurationSecs)
		}

		if bullet.Size <= 0 || bullet.Size > configs.Size_Max {
			return fmt.Errorf("invalid bullet[%d].Size (%d)", i, bullet.Size)
		}

		for j, v := range bullet.Color {
			if v < 0 || v > 1 {
				return fmt.Errorf("bullet[%d].Color[%d] (%f)", i, j, v)
			}
		}

		for j, v := range bullet.Position {
			if v > configs.Position_Max {
				return fmt.Errorf("bullet[%d].Position[%d] (%d)", i, j, v)
			}
		}
	}
	return nil
}

func CloseDrawingFile() {
	file.Sync()
	if err = file.Close(); err != nil {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			err = file.Close()
			if err != nil {
				pkg.WarnHandle(err, fmt.Sprintf("Closing the file not properly, try %d times...", i))
				continue
			}
		}
	}
	pkg.WarnHandle(err, "Failed on closing the file, please manually check the drawing file status")
	fmt.Println("CloseDrawingFile function ended!")
}

func WriteToDrawingFile(bullets *[]configs.Bullet) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	for _, v := range *bullets {
		fmt.Fprintf(file, "%v\n", v)
	}
	log.Printf("Wrote %d bullet(s) to the file!\n", len(*bullets))
}

func startDrawing() {
	log.Println("Drawing to OpenGL...")
}
