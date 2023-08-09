package main

import (
	"GRPC_TEST/service"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"os"
)

/*
证书分为根证书、服务器证书、客户端证书。根证书文件（ca.crt）和根证书对应的私钥文件（ca.key）由 CA（证书授权中心，国际认可）生成和保管。
那么服务器如何获得证书呢？向 CA 申请！步骤如下：
服务器生成自己的公钥（server.pub）和私钥（server.key）。后续通信过程中，客户端使用该公钥加密通信数据，服务端使用对应的私钥解密接收到的客户端的数据；
服务器使用公钥生成请求文件（server.req），请求文件中包含服务器的相关信息，比如域名、公钥、组织机构等；
服务器将 server.req 发送给 CA。CA 验证服务器合法后，使用 ca.key 和 server.req 生成证书文件（server.crt）——使用私钥生成证书的签名数据；
CA 将证书文件（server.crt）发送给服务器。
由于ca.key 和 ca.crt 是一对，ca.crt 文件中包含公钥，因此 ca.crt 可以验证 server.crt是否合法——使用公钥验证证书的签名。
*/

// key： 服务器上的私钥文件，用于对发送给客户端数据的加密，以及对从客户端接收到数据的解密。
// csr： 证书签名请求文件，用于提交给证书颁发机构（CA）对证书签名。
// crt： 由证书颁发机构（CA）签名后的证书，或者是开发者自签名的证书，包含证书持有人的信息，持有人的公钥，以及签署者的签名等信息。
// pem： 是基于Base64编码的证书格式，扩展名包括PEM、CRT和CER
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

	//Token,实现拦截器
	interceptor := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		err1 := AuthToken(ctx)
		if err1 != nil {
			return nil, err1
		}
		return handler(ctx, req)
	}
	fmt.Printf("Test")
	rpcServer := grpc.NewServer(grpc.Creds(newTLS), grpc.UnaryInterceptor(interceptor))
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

// AuthToken
// 验证Token是否合法
func AuthToken(c context.Context) error {
	//从上下文中读取数据
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return fmt.Errorf("can not found")
	}
	var user string
	var password string
	if val, ok1 := md["user"]; ok1 {
		user = val[0]
	}
	if val, ok1 := md["password"]; ok1 {
		password = val[0]
	}

	if user != "username" || password != "password" {
		return fmt.Errorf("username or password wrong")
	}
	return nil
}
