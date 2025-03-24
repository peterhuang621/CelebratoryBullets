package mylog

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"server/registry"
)

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	log.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	log.SetFlags(0)
	log.SetOutput(&clientLogger{url: serviceURL})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed to send log message. Service responded with code: %v", res.StatusCode)
	}
	return len(data), nil
}
