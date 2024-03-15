package main

import (
	// 导入自己项目的package
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
	host, port := "localhost", "6000"

	serviceAddress := fmt.Sprintf("http://%s:%s", host, port) 

	r := registry.Registeration {
		ServiceName: "Grade Service",
		ServiceURL: serviceAddress,
		RequiredService: []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress + "/services",
	}
	ctx, err := service.Start(
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

	if logProvider, err := registry.GetProvider(registry.LogService); err != nil {
		fmt.Printf("Logging service at: %v", logProvider)
		log.SetClientLogger(logProvider,r.ServiceName)
	}

	<- ctx.Done()
	fmt.Println("ShutDown")
}