package svc

import (
	"vehicle/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ServiceContext struct {
	Config config.Config
	MQTTClient mqtt.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MQTT 客户端
	opts := mqtt.NewClientOptions().AddBroker(c.MQTT.Broker)
	opts.SetClientID("vehicle_rpc_service")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &ServiceContext{
		Config:     c,
		MQTTClient: client,
	}
}
