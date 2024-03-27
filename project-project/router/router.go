package router

import (
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"log"
	"net"
	"strconv"
	"test.com/project-common/discovery"
	common "test.com/project-common/is_dev"
	"test.com/project-common/logs"
	"test.com/project-common/port"
	"test.com/project-grpc/account"
	"test.com/project-grpc/auth"
	"test.com/project-grpc/department"
	"test.com/project-grpc/menu"
	"test.com/project-grpc/project"
	"test.com/project-grpc/task"
	"test.com/project-project/config"
	"test.com/project-project/internal/interceptor"
	"test.com/project-project/internal/rpc"
	account_service_v1 "test.com/project-project/pkg/service/account.service.v1"
	auth_service_v1 "test.com/project-project/pkg/service/auth.service.v1"
	department_service_v1 "test.com/project-project/pkg/service/department.service.v1"
	menu_service_v1 "test.com/project-project/pkg/service/menu.service.v1"
	project_service_v1 "test.com/project-project/pkg/service/project.service.v1"
	task_service_v1 "test.com/project-project/pkg/service/task.service.v1"
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
	freePort, _ := port.FinAvailablePort(8000, 10000)
	addr := "127.0.0.1:" + strconv.Itoa(freePort)
	if common.IsDev {
		config.C.GC.Addr = addr
		config.C.GC.EtcdAddr = addr
	}

	c := gRPCConfig{
		Addr: config.C.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			project.RegisterProjectServiceServer(g, project_service_v1.New())
			task.RegisterTaskServiceServer(g, task_service_v1.New())
			account.RegisterAccountServiceServer(g, account_service_v1.New())
			department.RegisterDepartmentServiceServer(g, department_service_v1.New())
			auth.RegisterAuthServiceServer(g, auth_service_v1.New())
			menu.RegisterMenuServiceServer(g, menu_service_v1.New())
		}}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			//otelgrpc.UnaryServerInterceptor(),
			interceptor.New().CacheInterceptor(),
		)),
	)
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

func InitUserRpc() {
	rpc.InitRpcUserClient()
}
