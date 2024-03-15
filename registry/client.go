package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registeration) error {
	serviceUpdateURL,err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)

	if err != nil {
		return err
	}

	res, err := http.Post(ServiceURL, "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("注册失败 %v", res.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct {

}

func (shu *serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	prov.Update(p)
}

func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServiceURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to deregister %v", res.StatusCode)
	}
	return nil
}


type providers struct {
	service map[ServiceName][]string
	mutex *sync.RWMutex
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, pathcEntry := range pat.Added {
		if _, ok := p.service[pathcEntry.Name]; !ok {
			p.service[pathcEntry.Name] = make([]string, 0)
		}
		p.service[pathcEntry.Name] = append(p.service[pathcEntry.Name], pathcEntry.URL)
	}

	for _, pathcEntry := range pat.Removed {
		if providerURLs, ok := p.service[pathcEntry.Name]; !ok {
			for i := range providerURLs {
				if providerURLs[i] == pathcEntry.URL {
					p.service[pathcEntry.Name] = append(
						providerURLs[:i], 
						providerURLs[i+1:]...
					)
				}
			}
		}
	}
}

func (p providers) get(name ServiceName) (string,error) {
	providers,ok := p.service[name]
	if !ok {
		return "", fmt.Errorf("service %v not found", name)
	}
	idx := int(rand.Float32() * float32(len(providers)))
	return providers[idx], nil
}

func GetProvider(name ServiceName) (string,error) {
	return prov.get(name)
}

var prov = providers {
	service: make(map[ServiceName][]string),
	mutex: new(sync.RWMutex),
}