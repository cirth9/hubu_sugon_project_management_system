package main

import (
	"GRPC_TEST/client/auth"
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

var (
	username string //用户名
	password string //密码
)

func main() {
	//file1, err := credentials.NewClientTLSFromFile("client/cert_client/server.pem", "cirno.com")
	//if err != nil {
	//	log.Fatal("证书错误 ", err)
	//}
	_, err := fmt.Scan(&username, &password)
	if err != nil {
		log.Fatal("Input error")
	}
	user := &auth.Authentication{
		User:     username,
		Password: password,
	} //username 和 password都错了，预计输出password or username wrong

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
	newTLS := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{file},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ServerName: "*.cirnoblog.cn",
		RootCAs:    certPool,
	})

	//Dial第一个参数用于处理双向认证，第二个参数用于检测token是否合法，在此处指user和password是否合法
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(newTLS), grpc.WithPerRPCCredentials(user))
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
