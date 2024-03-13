package main

import (
	// 导入自己项目的package
	"Plt/log"
	"Plt/registry"
	"Plt/service"
	"context"
	"fmt"
	stlog "log"
)

func main() {
	// 运行Log程序
	log.Run("./destributed.log")

	// 注册Log服务
	host, port := "localhost", "4000"

	serviceAddress := fmt.Sprintf("http://%s:%s", host, port) 

	r := registry.Registeration {
		ServiceName: "Log Service",
		ServiceURL: serviceAddress,
	}
	ctx,err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandles,
	)

	// ShutDown
	if err != nil {
		stlog.Fatalln(err)
	}
	<- ctx.Done()
	fmt.Println("ShutDown")
}