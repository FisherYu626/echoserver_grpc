package main

import (
	"context"
	"fmt"
	"grpc/service"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	cred, err2 := credentials.NewClientTLSFromFile("../cert/server.pem", "*.fisher.com")
	if err2 != nil {
		log.Fatal("证书获取失败")
	}

	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(cred))

	if err != nil {
		log.Fatal("服务端连接出错，连接不上", err)
	}

	defer conn.Close()

	prodClient := service.NewProdServiceClient(conn)

	request := &service.ProductRequest{
		ProId: 12,
	}

	stockResponse, err := prodClient.GetProdStock(context.Background(), request)

	if err != nil {
		log.Fatal("查询库存出错", err)
	}

	fmt.Println("查询成功", stockResponse.ProdStock)

}
