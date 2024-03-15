package log

import (
	"io"
	stlog "log"
	"net/http"
	"os"
)


var log* stlog.Logger

type fileLog string

// 实现了这个方法的东西可以当 io.Writer 用 ，{ Write:func() } => io.Writer
func (fl fileLog) Write(data []byte) (int,error) {
	f,err := os.OpenFile(string(fl),os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0600)
	if err != nil {
		return 0,err
	}
	defer f.Close()
	return f.Write(data)
}

// cmd will call this func
func Run(destination string) {
	// 转换为这个自己的类型，就相当于重写了这个Write的方法
	log = stlog.New(fileLog(destination),"[go] - ",stlog.LstdFlags)
}

func write(message string) {
	log.Printf("%v\n",message)
}

// register
func RegisterHandles() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method{
		case http.MethodPost:
			msg,err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}
