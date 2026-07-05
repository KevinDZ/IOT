package config

import (
	"time"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MQTT struct {
		Broker         string
		StatusTopic    string
		ControlTopic   string
		ClientID       string
		CACert         string
		ClientCert     string
		ClientKey      string
		HeartbeatTopic string
		QoS            byte
	}
	Vehicle struct {
		ID string
	}
	LinuxPath         string
	WindowsPath       string
	Path              string
	HeartbeatInterval time.Duration
}
