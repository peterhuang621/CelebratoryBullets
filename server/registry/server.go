package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const ServerPort = ":3000"
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registration []Registration
	mutex        *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registration = append(r.registration, reg)
	r.mutex.Unlock()
	return nil
}

var reg = registry{
	registration: make([]Registration, 0),
	mutex:        new(sync.Mutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(ctx *gin.Context) {
	log.Println("Request received")
	switch ctx.Request.Method {
	case http.MethodPost:
		dec := json.NewDecoder(ctx.Request.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}
		log.Printf("Adding service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}
	default:
		ctx.JSON(http.StatusMethodNotAllowed, nil)
		return
	}
}
