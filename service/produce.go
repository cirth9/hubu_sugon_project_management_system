package service

import (
	"context"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/anypb"
)

type ProduceService struct {
}

func (prod *ProduceService) mustEmbedUnimplementedProdServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (prod *ProduceService) GetProdStockById(id int32) int32 {
	return id
}

func (prod *ProduceService) GetProdStock(ctx context.Context, in *Product_Request) (*Produce_Response, error) {
	stock := prod.GetProdStockById(in.ProdId)
	test := Test{
		Test: "test proto.Any",
	}
	anyVal, err := anypb.New(&test)
	if err != nil {
		return nil, err
	}

	return &Produce_Response{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		ProdStock:     stock,
		Data:          anyVal,
	}, nil
}
