package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const ServicePort = "3000"
const ServiceURL = "http://localhost" + ":" + ServicePort + "/services"

type registry struct {
	registerations []Registeration
	mutex	*sync.Mutex
}

func (r *registry) add(reg Registeration) error {
	r.mutex.Lock()
	r.registerations = append(r.registerations, reg)
	r.mutex.Unlock()
	return nil
}

func(r *registry) remove(url string) error {
	for i := range reg.registerations {
		if reg.registerations[i].ServiceURL == url {
			r.mutex.Lock()
			r.registerations = append(r.registerations[:i],r.registerations[i+1:]... )
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("找不到服务 %v",url)
}

// already init
var reg = registry {
	registerations: make([]Registeration, 0),
	mutex: new(sync.Mutex),
}

type RegistryService struct {}

func (s RegistryService) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registeration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("添加服务 %v 在 %v", r.ServiceName,r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
	case http.MethodDelete:
		payload,err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("Removing %v", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}