package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"fisher.com/grpc/client/auth"
	"fisher.com/grpc/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	//单向认证
	// cred, err2 := credentials.NewClientTLSFromFile("../cert/server.pem", "*.fisher.com")
	// if err2 != nil {
	// 	log.Fatal("证书获取失败")
	// }

	//双向认证
	cert, err := tls.LoadX509KeyPair("../cert/client.pem", "../cert/client.key")
	if err != nil {
		log.Fatal("证书获取错误", err)
	}

	//创建一个新的 空的 certpool
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../cert/ca.crt")
	if err != nil {
		log.Fatal("ca证书读取错误", err)
	}
	//尝试解析所传入的PEM编码的证书 如果解析成功会将其加到CertPool中 便于后面使用
	certPool.AppendCertsFromPEM(ca)
	//构建基于TLS的TransportCredentials选项
	creds := credentials.NewTLS(&tls.Config{
		//设置证书链
		Certificates: []tls.Certificate{cert},
		//要求必须校验客户端的证书 可以根据实际情况去选用以下参数
		ServerName: "*.fisher.com",
		RootCAs:    certPool,
	})

	//jwt方式或者Oauth方式 都可
	token := &auth.Authentication{
		User:   "admin",
		Passwd: "admin",
	}

	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(token))

	if err != nil {
		log.Fatal("服务端连接出错，连接不上", err)
	}

	defer conn.Close()

	prodClient := service.NewProdServiceClient(conn)

	// request := &service.ProductRequest{
	// 	ProId: 12,
	// }

	// stockResponse, err := prodClient.GetProdStock(context.Background(), request)

	// if err != nil {
	// 	log.Fatal("查询库存出错", err)
	// }

	// fmt.Println("查询成功", stockResponse.ProdStock)

	//客户端流
	// stream, err2 := prodClient.UpdateProdStockClientStream(context.Background())
	// if err2 != nil {
	// 	log.Fatal("获取流出错")
	// }
	// rsp := make(chan struct{}, 1)

	// go prodRequest(stream, rsp)

	// select {
	// case <-rsp:
	// 	recv, err3 := stream.CloseAndRecv()
	// 	if err3 != nil {
	// 		log.Fatal(err3)
	// 	}
	// 	stock := recv.ProdStock
	// 	fmt.Println("客户端收到响应", stock)
	// }

	//服务端流
	// request := &service.ProductRequest{
	// 	ProId: 12,
	// }
	// stream, err2 := prodClient.GetProdStockServerStream(context.Background(), request)
	// if err2 != nil {
	// 	log.Fatal("流出错")
	// }

	// for {
	// 	recv, err := stream.Recv()
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			fmt.Println("客户端数据接收完毕")
	// 			err := stream.CloseSend()
	// 			if err != nil {
	// 				log.Fatal(err)
	// 			}
	// 			break
	// 		}
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println("客户端收到的流", recv.ProdStock)
	// }

	//双向流
	stream, err2 := prodClient.GetProdStockDoubleStream(context.Background())
	if err2 != nil {
		log.Fatal(err2)
	}

	for {
		request := &service.ProductRequest{
			ProId: 12,
		}
		err3 := stream.Send(request)
		if err3 != nil {
			log.Fatal(err3)
		}

		time.Sleep(time.Second)

		recv, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("客户端收到的流信息", recv.ProdStock)

	}
}

func prodRequest(stream service.ProdService_UpdateProdStockClientStreamClient, rsp chan struct{}) {
	count := 0
	for {
		request := &service.ProductRequest{
			ProId: 12,
		}

		err := stream.Send(request)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)

		count++

		if count > 10 {
			rsp <- struct{}{}
			break
		}
	}
}
