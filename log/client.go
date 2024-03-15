package log

import (
	"bytes"
	"fmt"
	"go-distuibuted/registry"
	stlog "log"
	"net/http"
)




func SetClientLogger(ServiceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ",clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: ServiceURL})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int,error) {
	b := bytes.NewBuffer([]byte(data))
	res, err := http.Post(cl.url+"/log", "application/json", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("fail to send log message: %d", data)
	}
	return len(data), nil
}