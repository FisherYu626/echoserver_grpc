package main

import (
	"fmt"
	"grpc/service"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	//无认证 http2 https

	cred, err2 := credentials.NewServerTLSFromFile("cert/server.pem", "cert/server.key")
	if err2 != nil {
		log.Fatal("读取证书错误")
	}

	rpcServer := grpc.NewServer(grpc.Creds(cred))

	service.RegisterProdServiceServer(rpcServer, service.ProductService)

	listener, err := net.Listen("tcp", ":8002")

	if err != nil {
		log.Fatal("启动监听出错", err)
	}

	err = rpcServer.Serve(listener)

	if err != nil {
		log.Fatal("启动服务出错", err)
	}

	fmt.Println("启动服务成功")
}
