package svc

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
	"vehicle/internal/config"
	"vehicle/pb/vehicle"
	"vehicle/vehicleservice"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	MQTTClient mqtt.Client
	Ctx        context.Context
	Cancel     context.CancelFunc
	SigChan    chan os.Signal
}

func NewServiceContext(c config.Config, wg *sync.WaitGroup) *ServiceContext {
	// 初始化 MQTT 客户端
	opts := mqtt.NewClientOptions().AddBroker(c.MQTT.Broker)
	opts.SetClientID(c.MQTT.ClientID)
	opts.SetAutoReconnect(true)

	// 🚨 核心：配置遗嘱消息 (Will Message)
	// 当 EMQX 检测到该客户端异常掉线（如突然断电、拔网线）时，
	// 会自动向 "/vehicle/status/upload" 发布这条离线消息。
	willPayload := `{"status":"offline", "reason":"unexpected_disconnect"}`
	opts.SetWill(
		"/vehicle/status/upload", // 遗嘱发布的 Topic
		willPayload,              // 遗嘱消息体 (建议用 JSON 格式)
		1,                        // QoS 级别 (1 保证至少到达一次)
		true,                     // Retain 标志 (true 表示服务端会保留这条离线状态，新订阅者也能看到)
	)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("⚠️ [Fallback] 收到未匹配的消息: Topic=%s, Payload=%s", msg.Topic(), string(msg.Payload()))
	})

	// 2. 设置连接成功后的回调（在这里进行订阅，确保不会丢失消息）
	opts.OnConnect = func(client mqtt.Client) { // 确保这个 Topic 与手机端发布的完全一致
		token := client.Subscribe(c.MQTT.ControlTopic, 1, handleCommand) // QoS 设为 1
		token.Wait()
		if token.Error() != nil {
			log.Printf("❌ 订阅 Topic [%s] 失败: %v\n", c.MQTT.ControlTopic, token.Error())
		} else {
			log.Printf("✅ 成功订阅 Topic: %s\n", c.MQTT.ControlTopic)
		}
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Printf("❌ 连接丢失: %v, 错误类型: %T", err, err)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// 或者使用带超时的等待
	if !token.WaitTimeout(5 * time.Second) {
		log.Fatal("❌ 连接超时")
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())

	s := &ServiceContext{
		Config:     c,
		MQTTClient: client,
		Ctx:        ctx,
		Cancel:     cancel,
		SigChan:    make(chan os.Signal),
	}

	s.initMQTT(wg)

	return s
}

func (s *ServiceContext) initMQTT(wg *sync.WaitGroup) {
	log.Println("✅ MQTT 连接成功，正在启动后台协程...")

	// 启动协程 1：监听控制指令
	go s.listenCommands(wg)

	// 启动协程 2：定时上传车辆状态
	go s.uploadStatus(wg)
}

func (sc *ServiceContext) listenCommands(wg *sync.WaitGroup) {

	wg.Add(1)
	log.Println("🚗 车载端已启动，等待手机端指令...")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop() // 确保协程退出时释放定时器资源

	for {
		select {
		case <-sc.Ctx.Done(): // 等待上下文取消信号
			log.Println("🚗 车载端已安全退出")
			wg.Done()
			return
		case <-ticker.C:
			token := sc.MQTTClient.Subscribe(sc.Config.MQTT.ControlTopic, 1, handleCommand)

			if !token.WaitTimeout(5 * time.Second) {
				// 1. 使用带超时的等待，防止网络卡死导致无法响应取消信号
				log.Printf("⚠️ 订阅 Topic [%s] 超时，将在下个周期重试\n", sc.Config.MQTT.ControlTopic)
				continue // 超时不退出，继续重试
			}
			if token.Error() != nil {
				// 2. 订阅失败不要 return，改为 continue 等待下一次重试
				log.Printf("❌ 订阅 Topic [%s] 失败: %v，将在下个周期重试\n", sc.Config.MQTT.ControlTopic, token.Error())
				continue
			}
			// 3. 订阅成功后，打印成功信息
			log.Printf("✅ 成功订阅 Topic: %s\n", sc.Config.MQTT.ControlTopic)
		}
	}
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

func (sc *ServiceContext) uploadStatus(wg *sync.WaitGroup) {
	wg.Add(1)

	// 定时上传车辆状态
	battery := 85.0 // todo 实际项目中这里读取 CAN 总线或传感器数据

	ticker := time.NewTicker(5 * time.Second) // 每5秒上报一次
	defer ticker.Stop()

	for {
		select {
		case <-sc.Ctx.Done():
			log.Println("🛑 [状态上传] 收到退出信号，协程安全退出")
			wg.Done()
			return
		case <-ticker.C:
			// 💡 建议：正常上报时，明确带上 "online" 状态
			isOnline := true
			speed := 65.5
			status := vehicle.VehicleStatusResp{
				VehicleID:   sc.Config.Vehicle.ID,
				Speed:       speed, // todo 实际项目中这里读取 CAN 总线或传感器数据
				Battery:     battery,
				IsOnline:    isOnline,
				LastUpdated: float64(time.Now().Unix()),
			}

			payload, err := json.Marshal(status)
			if err != nil {
				log.Printf("❌ 序列化状态失败: %v\n", err)
				continue
			}
			// 发布到专属的状态 Topic
			token := sc.MQTTClient.Publish(sc.Config.MQTT.StatusTopic, 1, false, payload)
			go func(t mqtt.Token) {
				token.Wait()
				if token.Error() != nil {
					log.Printf("❌ 发布状态失败: %v\n", token.Error())
					return
				}
				log.Printf("🚗 定时更新汽车状态信息：%v\n", status)
			}(token)
			// todo 测试用
			if battery < 0 {
				battery = 100
			}
			battery--
		}
	}
}

func (sc *ServiceContext) ShutDown(service *zrpc.RpcServer, quit *chan os.Signal, shutdownDone *chan struct{}, wg *sync.WaitGroup, once *sync.Once) {
	go func() {
		<-*quit
		log.Println("🛑 收到系统退出信号，开始清理...")

		// 使用 sync.Once 确保清理逻辑只执行一次
		once.Do(func() {
			// 触发全局广播
			sc.Cancel()

			// 等待所有协程安全退出
			wg.Wait()
			log.Println("✅ 所有协程已安全退出")

			// 停止 RPC 服务
			log.Println("🛑 开始停止 RPC 服务...")
			service.Stop()
			close(*shutdownDone)
			log.Println("✅ RPC 服务已停止")
		})
	}()
}
