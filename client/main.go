package main

import (
	"GRPC_TEST/service"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.Dial(":8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	rpcClient := service.NewProdServiceClient(conn)
	resp, err := rpcClient.GetProdStock(context.Background(), &service.Product_Request{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", resp)
}
