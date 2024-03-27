package main

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"log"
	"strconv"
	srv "test.com/project-common"
	"test.com/project-common/port"
	"test.com/project-user/config"
	"test.com/project-user/router"
	"test.com/project-user/tracing"

	common "test.com/project-common/is_dev"
)

func main() {
	r := gin.Default()
	tp, tpErr := tracing.JaegerTraceProvider(config.C.JaegerConfig.Endpoints)
	if tpErr != nil {
		log.Fatal(tpErr)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	//路由
	router.InitRouter(r)
	//grpc服务注册
	gc := router.RegisterGrpc()
	//grpc服务注册到etcd
	router.RegisterEtcdServer()

	stop := func() {
		gc.Stop()
	}
	if !common.IsDev {
		srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
	} else {
		freePort, _ := port.FinAvailablePort(8000, 10000)
		addr := "127.0.0.1:" + strconv.Itoa(freePort)
		zap.S().Info("project-user srv address:", addr)
		srv.Run(r, config.C.SC.Name, addr, stop)
	}
}
