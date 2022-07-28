package service

import (
	context "context"
	"fmt"
	"io"
	"log"
	"time"
)

var ProductService = &productService{}

type productService struct {
}

func (p *productService) GetProdStock(conetxt context.Context, request *ProductRequest) (*ProductResponse, error) {
	//实现具体业务逻辑
	stock := p.GetStcokById(request.ProId)

	return &ProductResponse{ProdStock: stock}, nil
}

func (p *productService) GetStcokById(id int32) int32 {
	return id + 1
}

func (p *productService) mustEmbedUnimplementedProdServiceServer() {

}

func (p *productService) UpdateProdStockClientStream(stream ProdService_UpdateProdStockClientStreamServer) error {
	count := 0
	for {
		//源源不断接收客户端发来的信息
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		count++
		fmt.Println("服务端接收到的流", recv.ProId, count)
		if count > 10 {
			rsp := &ProductResponse{ProdStock: recv.ProId + 2}
			err := stream.SendAndClose(rsp)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (p *productService) GetProdStockServerStream(request *ProductRequest, stream ProdService_GetProdStockServerStreamServer) error {
	count := 0
	for {
		rsp := &ProductResponse{ProdStock: request.ProId + 2}

		err := stream.Send(rsp)
		if err != nil {
			return err
		}

		time.Sleep(time.Second)
		count++
		if count > 10 {
			return nil
		}
	}

	return nil
}

func (p *productService) GetProdStockDoubleStream(stream ProdService_GetProdStockDoubleStreamServer) error {

	for {
		recv, err := stream.Recv()
		if err != nil {
			return nil
		}
		fmt.Println("服务端收到客户端消息", recv.ProId)
		time.Sleep(time.Second)

		rsp := &ProductResponse{ProdStock: recv.ProId + 2}
		err2 := stream.Send(rsp)
		if err2 != nil {
			log.Fatal(err2)
		}

	}

	return nil
}
