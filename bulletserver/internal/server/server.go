package server

import (
	"fmt"
	"os"
	"time"

	"github.com/peterhuang621/CelebratoryBullets/bulletserver/configs"
	"github.com/peterhuang621/CelebratoryBullets/bulletserver/pkg"
)

var file *os.File
var err error

func init() {
	fileInfo, err := os.Stat(configs.GL_Drawingfile)
	if os.IsNotExist(err) {
		file, err = os.Create(configs.GL_Drawingfile)
		pkg.ErrorHandle(err, "Failed on creating the nonexisted file")
	} else {
		file, err = os.Open(configs.GL_Drawingfile)
		pkg.ErrorHandle(err, "Failed on opening the existed file")
		if fileInfo.Size() != 0 {
			startDrawing()
		}
	}
}

func CloseDrawingFile() {
	if err = file.Close(); err != nil {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			err = file.Close()
			if err != nil {
				pkg.ErrorHandle(err, fmt.Sprintf("Closing the file not properly, try %d times...", i))
				continue
			}
		}
	}
	pkg.ErrorHandle(err, "Failed on closing the file")
}

func WriteToDrawingFile(bullets []configs.Bullet) {
	for _, v := range bullets {
		fmt.Fprintf(file, "%v\n", v)
	}
}

func startDrawing() {

}
