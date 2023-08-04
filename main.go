package main

import (
	"GRPC_TEST/service"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	fmt.Printf("Test")
	rpcServer := grpc.NewServer()
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
