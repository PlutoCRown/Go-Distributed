package main

import (
	// 导入自己项目的package
	"context"
	"fmt"
	"go-distuibuted/grades"
	"go-distuibuted/log"
	"go-distuibuted/registry"
	"go-distuibuted/service"
	stlog "log"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registeration{
		ServiceName:      "Grade Service",
		ServiceURL:       serviceAddress,
		RequiredService:  []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress + "/services",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers,
	)

	// ShutDown
	if err != nil {
		stlog.Fatalln(err)
	}
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("找到日志服务: %v", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	} else {
		fmt.Printf("错误: %v", err)
	}

	<-ctx.Done()
	fmt.Println("ShutDown")
}
