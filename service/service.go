package service

import (
	"context"
	"fmt"
	"go-distuibuted/registry"
	"log"
	"net/http"
)

func Start(
	ctx context.Context,
	host,
	port string,
	reg registry.Registeration,
	registerHandlesFunc func(),
) (context.Context, error) {

	registerHandlesFunc()
	ctx = startService(ctx, reg.ServiceName, host, port)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, serviceName registry.ServiceName, host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var srv http.Server
	srv.Addr = ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		fmt.Println("应该注销服务1")
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		fmt.Printf("%v 服务已启动，按任意键暂停\n", serviceName)
		var a string
		fmt.Scanln(&a)
		fmt.Println("应该注销服务2")
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
