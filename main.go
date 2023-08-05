package main

import (
	"GRPC_TEST/service"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

func main() {
	//添加证书
	file1, err := credentials.NewServerTLSFromFile("cert/server.pem", "cert/server.key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Test")
	rpcServer := grpc.NewServer(grpc.Creds(file1))

	service.RegisterProdServiceServer(rpcServer, &service.ProduceService{})
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	defer func(listener net.Listener) {
		err = listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}

}
