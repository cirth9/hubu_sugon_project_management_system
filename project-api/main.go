package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log"
	"net/http"
	_ "test.com/project-api/api"
	"test.com/project-api/api/midd"
	"test.com/project-api/config"
	"test.com/project-api/router"
	"test.com/project-api/tracing"
	srv "test.com/project-common"
)

func main() {
	r := gin.Default()
	tp, tpErr := tracing.JaegerTraceProvider(config.C.JaegerConfig.Endpoints)
	if tpErr != nil {
		log.Fatal(tpErr)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	r.Use(midd.RequestLog())
	r.Use(otelgin.Middleware("project-api"))

	r.StaticFS("/upload", http.Dir("upload"))
	//路由
	router.InitRouter(r)
	//开启pprof 默认的访问路径是/debug/pprof
	//pprof.Register(r)
	////测试代码
	//r.GET("/mem", func(c *gin.Context) {
	//	// 业务代码运行
	//	outCh := make(chan int)
	//	// 每秒起10个goroutine，goroutine会阻塞，不释放内存
	//	tick := time.Tick(time.Second / 10)
	//	i := 0
	//	for range tick {
	//		i++
	//		fmt.Println(i)
	//		alloc1(outCh) // 不停的有goruntine因为outCh堵塞，无法释放
	//	}
	//})
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, nil)
}

// 一个外层函数
func alloc1(outCh chan<- int) {
	go alloc2(outCh)
}

// 一个内层函数
func alloc2(outCh chan<- int) {
	func() {
		defer fmt.Println("alloc-fm exit")
		// 分配内存，假用一下
		buf := make([]byte, 1024*1024*10)
		_ = len(buf)
		fmt.Println("alloc done")

		outCh <- 0
		//return
	}()
}
