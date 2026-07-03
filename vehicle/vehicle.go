package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"log"
	"vehicle/internal/config"
	"vehicle/internal/server"
	"vehicle/internal/svc"
	"vehicle/pb/vehicle"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/vehicle.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	var wg sync.WaitGroup

	conf.MustLoad(*configFile, &c)

	log.Println(c.MQTT.CACert, c.MQTT.ClientCert, c.MQTT.ClientKey)
	tls := config.NewTLS(c.Path, c.MQTT.CACert, c.MQTT.ClientCert, c.MQTT.ClientKey)
	log.Println(tls)

	ctx := svc.NewServiceContext(c, tls, &wg)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		vehicle.RegisterVehicleServiceServer(grpcServer, server.NewVehicleServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// defer s.Stop()

	// 3. 启动独立的协程监听系统信号（Ctrl+C 或 K8s 的 SIGTERM）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var once sync.Once
	var shutdownDone chan struct{}
	shutdownDone = make(chan struct{})

	// 优雅关闭
	ctx.ShutDown(s, &quit, &shutdownDone, &wg, &once)

	// 2. 最后再启动 RPC 服务（主协程在此阻塞）
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	go s.Start()

	<-shutdownDone
	log.Println("👋 进程正常退出")
}
