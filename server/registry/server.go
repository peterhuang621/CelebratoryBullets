package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const ServerPort = ":3000"
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

func (r *registry) remove(url string) error {
	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			reg.mutex.Lock()
			reg.registrations = append(reg.registrations[:i], reg.registrations[i+1:]...)
			reg.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("Service at URL %s not found", url)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
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
	case http.MethodDelete:
		payload, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
		url := string(payload)
		log.Printf("Removing service at URL: %s", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}
	default:
		ctx.JSON(http.StatusMethodNotAllowed, nil)
		return
	}
}
