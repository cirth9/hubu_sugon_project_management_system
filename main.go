package main

import (
	"GRPC_TEST/service"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
)

func main() {
	//添加证书
	//file1, err := credentials.NewServerTLSFromFile("cert/server.pem", "cert/server.key")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//证书认证，双向认证
	//从证书相关文件中读取和解析信息，得到证书的公钥秘钥对
	file, err := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
	if err != nil {
		log.Fatal(err)
	}
	//创建certPool（证书池）
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("cert/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	//尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	// 构建基于 TLS 的 TransportCredentials 选项
	newTLS := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{file},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ClientAuth: tls.RequireAndVerifyClientCert,
		// 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
		ClientCAs: certPool,
	})
	fmt.Printf("Test")
	rpcServer := grpc.NewServer(grpc.Creds(newTLS))

	service.RegisterProdServiceServer(rpcServer, &service.ProduceService{})
	listener, err := net.Listen("tcp", ":8002")
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
