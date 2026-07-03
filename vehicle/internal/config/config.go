package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MQTT struct {
		Broker       string
		StatusTopic  string
		ControlTopic string
		ClientID     string
		CACert       string
		ClientCert   string
		ClientKey    string
	}
	Vehicle struct {
		ID string
	}
	Path string
}
