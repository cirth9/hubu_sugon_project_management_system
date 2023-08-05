package main

import (
	"GRPC_TEST/service"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func main() {
	file1, err := credentials.NewClientTLSFromFile("client/cert_client/server.pem", "cirno.com")
	if err != nil {
		log.Fatal("证书错误 ", err)
	}

	conn, err := grpc.Dial(":8000", grpc.WithTransportCredentials(file1))
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
