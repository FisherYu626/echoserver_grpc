package service

import context "context"

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
