package registry

import (
	"bytes"
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
	mutex          *sync.RWMutex
}

func (r *registry) add(reg Registeration) error {
	r.mutex.Lock()
	r.registerations = append(r.registerations, reg)
	r.mutex.Unlock()
	err := r.sendRequiredServices(reg)
	r.notity(patch{
		Added: []patchEntry{
			patchEntry{
				Name: reg.ServiceName,
				URL:  reg.ServiceURL,
			},
		},
	})
	return err
}

// 反向依赖处理：上下线时查找依赖的服务
func (r registry) notity(fullPatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, reg := range r.registerations {
		go func(reg Registeration) {
			for _, reqService := range reg.RequiredService {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false
				for _, added := range fullPatch.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range fullPatch.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

// 服务上线时主动查找
func (r *registry) sendRequiredServices(reg Registeration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registerations {
		for _, reqService := range reg.RequiredService {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

// 发送更新包到更新接口
func (r registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) remove(url string) error {
	for i := range reg.registerations {
		if reg.registerations[i].ServiceURL == url {
			r.notity(patch{
				Added: []patchEntry{},
				Removed: []patchEntry{{
					Name: r.registerations[i].ServiceName,
					URL:  r.registerations[i].ServiceURL,
				}},
			})
			r.mutex.Lock()
			r.registerations = append(r.registerations[:i], r.registerations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("找不到服务 %v", url)
}

// already init
var reg = registry{
	registerations: make([]Registeration, 0),
	mutex:          new(sync.RWMutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("添加服务 %v 在 %v", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	case http.MethodDelete:
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("服务下线 %v", url)
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
