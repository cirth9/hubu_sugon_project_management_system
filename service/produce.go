package service

import (
	"context"
	"google.golang.org/protobuf/runtime/protoimpl"
)

type ProduceService struct {
}

func (prod *ProduceService) GetProdStock(ctx context.Context, in *Product_Request) (*Produce_Response, error) {
	return &Produce_Response{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		ProdStock:     18,
	}, nil
}
