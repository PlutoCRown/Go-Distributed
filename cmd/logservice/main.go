package main

import (
	"context"
	"fmt"

	"go-distuibuted/log"
	"go-distuibuted/registry"
	"go-distuibuted/service"
	stlog "log"
)

func main() {
	// 运行Log程序
	log.Run("./destributed.log")
	// 注册Log服务
	host, port := "localhost", "4000"

	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registeration{
		ServiceName:      "Log Service",
		ServiceURL:       serviceAddress,
		RequiredService:  make([]registry.ServiceName, 0),
		ServiceUpdateURL: serviceAddress + "/services",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)

	// ShutDown
	if err != nil {
		stlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("ShutDown")
}
