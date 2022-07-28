package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"fisher.com/grpc/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		fmt.Errorf("获取用户信息失败")
	}

	var user string
	var passwd string

	if val, ok := md["user"]; ok {
		user = val[0]
	}

	if val, ok := md["passwd"]; ok {
		passwd = val[0]
	}

	if user != "admin" || passwd != "admin" {
		return status.Errorf(codes.Unauthenticated, "token不合法")
	}

	return nil
}

func main() {
	//无认证 http2 https

	//单向认证
	// cred, err2 := credentials.NewServerTLSFromFile("cert/server.pem", "cert/server.key")
	// if err2 != nil {
	// 	log.Fatal("读取证书错误")
	// }

	//证书认证 双向认证
	//从证书相关文件读取和解析信息 得到证书公钥 私钥
	cert, err := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
	if err != nil {
		log.Fatal("证书获取错误", err)
	}

	//创建一个新的 空的 certpool
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("cert/ca.crt")
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
		ClientAuth: tls.RequireAndVerifyClientCert,
		//设置根证书的集合 校验方式用ClientAuth指定的方式
		ClientCAs: certPool,
	})

	//实现token认证 需要合法的用户名与密码
	//实现一个拦截器

	var authInterceptor grpc.UnaryServerInterceptor

	authInterceptor = func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		//拦截方法请求 验证token
		err = Auth(ctx)

		if err != nil {
			log.Fatal("身份认证失败")
		}

		//继续处理请求
		return handler(ctx, req)
	}

	rpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(authInterceptor))

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
