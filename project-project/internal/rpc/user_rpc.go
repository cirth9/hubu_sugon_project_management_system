package rpc

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"test.com/project-common/discovery"
	"test.com/project-common/logs"
	"test.com/project-grpc/user/login"
	"test.com/project-project/config"
)

var LoginServiceClient login.LoginServiceClient

func InitRpcUserClient() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	conn, err := grpc.Dial(
		"etcd:///user?version=1.0.0",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	LoginServiceClient = login.NewLoginServiceClient(conn)
}
