package main

import (
	"GRPC_TEST/service"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
)

func main() {
	//file1, err := credentials.NewClientTLSFromFile("client/cert_client/server.pem", "cirno.com")
	//if err != nil {
	//	log.Fatal("证书错误 ", err)
	//}
	file, err := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
	if err != nil {
		log.Fatal(1, err)
	}
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("cert/ca.crt")
	if err != nil {
		log.Fatal(2, err)
	}
	certPool.AppendCertsFromPEM(ca)
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{file},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ServerName: "*.cirnoblog.cn",
		RootCAs:    certPool,
	})

	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(3, err)
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			log.Fatal(4, err)
		}
	}(conn)

	rpcClient := service.NewProdServiceClient(conn)
	request := &service.Product_Request{
		ProdId: 123,
	}
	resp, err := rpcClient.GetProdStock(context.Background(), request)
	if err != nil {
		log.Fatal("response failed ", err)
	}
	fmt.Printf("%#v", resp)
}
