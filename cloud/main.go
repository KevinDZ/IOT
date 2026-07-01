package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// 极简的内存级 MQTT 消息路由中心
// 生产环境中，请直接部署 EMQX 或 VerneMQ

var (
	clients = make(map[net.Conn]string) // 存储已连接的客户端及其订阅的 Topic
	mu      sync.Mutex                  // 并发锁，保证多客户端连接时的数据安全
)

func main() {
	listener, err := net.Listen("tcp", ":1883")
	if err != nil {
		log.Fatalf("☁️ 云端中枢启动失败: %v", err)
	}
	defer listener.Close()

	log.Println("☁️ 云端中枢 (Cloud Broker) 已启动，监听端口: 1883")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("⚠️ 接受连接异常: %v", err)
			continue
		}
		// 为每个新连接开启一个独立的 Goroutine 进行处理
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		conn.Close()
		log.Printf("🔌 客户端断开连接: %s", conn.RemoteAddr())
	}()

	log.Printf("🔗 新客户端接入: %s", conn.RemoteAddr())
	reader := bufio.NewReader(conn)

	for {
		// 读取客户端发来的消息（以换行符为界）
		message, err := reader.ReadString('\n')
		if err != nil {
			return // 客户端断开或发生错误
		}

		message = strings.TrimSpace(message)
		parts := strings.SplitN(message, " ", 3)
		if len(parts) < 2 {
			continue
		}

		action := parts[0]
		switch action {
		case "SUBSCRIBE": // 处理订阅请求 (车载端使用)
			topic := parts[1]
			mu.Lock()
			clients[conn] = topic
			mu.Unlock()
			log.Printf("📥 客户端 %s 订阅了 Topic: %s", conn.RemoteAddr(), topic)

		case "PUBLISH": // 处理发布请求 (手机端API使用)
			if len(parts) < 3 {
				continue
			}
			topic := parts[1]
			payload := parts[2]
			routeMessage(topic, payload)
		}
	}
}

// 核心路由逻辑：将消息精准投递给订阅了对应 Topic 的客户端
func routeMessage(topic, payload string) {
	mu.Lock()
	defer mu.Unlock()

	log.Printf("🔄 云端收到消息 -> Topic: %s, Payload: %s", topic, payload)

	for conn, subTopic := range clients {
		// 简单的 Topic 匹配逻辑（生产环境需支持通配符如 + 和 #）
		if subTopic == topic {
			msg := fmt.Sprintf("MSG %s %s\n", topic, payload)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Printf("⚠️ 消息投递失败: %v", err)
			} else {
				log.Printf("✅ 消息已推送至客户端: %s", conn.RemoteAddr())
			}
		}
	}
}
