package mylog

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var serverlog *log.Logger

type fileLog string

func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

func Run(des string) {
	serverlog = log.New(fileLog(des), "go server: ", log.LstdFlags)
}

func RegisterHandlers(r *gin.Engine) {

	r.POST("/log", func(ctx *gin.Context) {
		msg, err := io.ReadAll(ctx.Request.Body)
		if err != nil || len(msg) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
		} else {
			write(string(msg))
			ctx.JSON(http.StatusOK, gin.H{
				"status": "message received",
			})
		}
	})

	r.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "method not allowed",
		})
	})

	// api.Any("", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusMethodNotAllowed, gin.H{
	// 		"error": "not allowed request",
	// 	})
	// })

}

func write(message string) {
	serverlog.Printf("%v\n", message)
}
