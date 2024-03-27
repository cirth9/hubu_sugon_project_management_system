package router

import (
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"strconv"
	"test.com/project-common/port"
	"test.com/project-grpc/user/user_basic"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"google.golang.org/grpc/resolver"
	"log"

	"net"
	"test.com/project-common/discovery"

	"test.com/project-common/logs"
	"test.com/project-grpc/user/login"

	"test.com/project-user/config"
	loginServiceV1 "test.com/project-user/pkg/service/login.service.v1"
	userBasicServiceV1 "test.com/project-user/pkg/service/user_basic.service.v1"

	common "test.com/project-common/is_dev"
)

// Router 接口
type Router interface {
	Route(r *gin.Engine)
}

type RegisterRouter struct {
}

func New() *RegisterRouter {
	return &RegisterRouter{}
}

func (*RegisterRouter) Route(ro Router, r *gin.Engine) {
	ro.Route(r)
}

var routers []Router

func InitRouter(r *gin.Engine) {
	//rg := New()
	//rg.Route(&user.RouterUser{}, r)
	for _, ro := range routers {
		ro.Route(r)
	}
}

func Register(ro ...Router) {
	routers = append(routers, ro...)
}

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGrpc() *grpc.Server {
	//0.0.0.0:8881
	freePort, _ := port.FinAvailablePort(8000, 10000)
	addr := "127.0.0.1:" + strconv.Itoa(freePort)
	if common.IsDev {
		config.C.GC.Addr = addr
		config.C.GC.EtcdAddr = addr
	}
	zap.S().Info("project-user grpc srv address:", config.C.GC.Addr)

	c := gRPCConfig{
		Addr: config.C.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			login.RegisterLoginServiceServer(g, loginServiceV1.New())
			user_basic.RegisterUserBasicServiceServer(g, userBasicServiceV1.New())
		}}
	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		otelgrpc.UnaryServerInterceptor(),
		//interceptor.New().CacheInterceptor(),
	)))
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	go func() {
		log.Printf("grpc server started as: %s \n", c.Addr)
		err = s.Serve(lis)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()
	return s
}

func RegisterEtcdServer() {
	etcdRegister := discovery.NewResolver(config.C.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)
	//服务地址:8881
	info := discovery.Server{
		Name:    config.C.GC.Name,
		Addr:    config.C.GC.EtcdAddr,
		Version: config.C.GC.Version,
		Weight:  config.C.GC.Weight,
	}
	zap.L().Info("register grpc addr: ", zap.String("addr", info.Addr))
	r := discovery.NewRegister(config.C.EtcdConfig.Addrs, logs.LG)
	_, err := r.Register(info, 2)
	if err != nil {
		log.Fatalln(err)
	}
}
