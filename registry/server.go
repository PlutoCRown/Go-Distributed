package registry

import (
	"encoding/json"
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

func(r *registry) remove(reg Registeration) error {
	r.mutex.Lock()
	// for( i in r.registerations) {

	// }
	// r.registerations = append(r.registerations[:i],r.registerations[i+1]...)
	r.mutex.Unlock()
	return  nil
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
		dec := json.NewDecoder(r.Body)
		var r Registeration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("移除服务 %v 在 %v", r.ServiceName,r.ServiceURL)
		err = reg.remove(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
	}
}