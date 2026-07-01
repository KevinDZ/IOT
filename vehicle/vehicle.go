package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"log"
	"os"
	"os/signal"
	"syscall"
	"vehicle/internal/config"
	"vehicle/internal/server"
	"vehicle/internal/svc"
	"vehicle/vehicle"
	"vehicle/vehicleservice"

	mqtt "github.com/eclipse/paho.mqtt.golang"

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
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	go func() {
		// 1. 配置 EMQX 连接选项
		opts := mqtt.NewClientOptions()
		opts.AddBroker("tcp://172.28.29.226:1883") // 替换为你的 EMQX 地址
		opts.SetClientID("vehicle-rpc-1001")
		opts.SetAutoReconnect(true)

		// 2. 设置连接成功后的回调（在这里进行订阅，确保不会丢失消息）
		opts.OnConnect = func(client mqtt.Client) {
			topic := "/vehicle/1001/control"                   // 确保这个 Topic 与手机端发布的完全一致
			token := client.Subscribe(topic, 1, handleCommand) // QoS 设为 1
			token.Wait()
			if token.Error() != nil {
				log.Printf("❌ 订阅 Topic [%s] 失败: %v\n", topic, token.Error())
			} else {
				log.Printf("✅ 成功订阅 Topic: %s\n", topic)
			}
		}

		// 3. 建立连接
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatalf("❌ 连接 EMQX 失败: %v", token.Error())
		}

		log.Println("🚗 车载端已启动，等待手机端指令...")

		// 保持程序运行，直到收到中断信号
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		client.Disconnect(250)
		log.Println("🚗 车载端已安全退出")
	}()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		vehicle.RegisterVehicleServer(grpcServer, server.NewVehicleServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

// 4. 消息处理回调函数
func handleCommand(client mqtt.Client, msg mqtt.Message) {
	log.Printf("📥 收到原始消息: %s\n", msg.Payload())

	var cmd vehicleservice.ControlReq
	// 解析 JSON
	if err := json.Unmarshal(msg.Payload(), &cmd); err != nil {
		log.Printf("❌ JSON 解析失败: %v\n", err)
		return
	}

	// 安全校验：检查 vehicleID 是否缺失（解决之前的报错）
	if cmd.VehicleID == "" {
		log.Println("❌ 错误: field \"vehicleID\" is not set")
		return
	}

	// 5. 根据指令执行具体动作
	log.Printf("🔧 收到控制指令: VehicleID=%s, Action=%s\n", cmd.VehicleID, cmd.Action)
	switch cmd.Action {
	case "unlock_door":
		log.Println("🔓 执行动作: 解锁车门")
		// TODO: 调用底层硬件接口或 CAN 总线解锁
	case "lock_door":
		log.Println("🔒 执行动作: 锁定车门")
	default:
		log.Printf("⚠️ 未知指令: %s\n", cmd.Action)
	}
}
