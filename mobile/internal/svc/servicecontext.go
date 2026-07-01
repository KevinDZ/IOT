// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"mobile/internal/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ServiceContext struct {
	Config     config.Config
	MQTTClient mqtt.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MQTT 客户端
	opts := mqtt.NewClientOptions().AddBroker(c.MQTT.Broker)
	opts.SetClientID("mobile_api_go_zero")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return &ServiceContext{
		Config:     c,
		MQTTClient: client,
	}
}
