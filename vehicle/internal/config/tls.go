package config

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

func NewTLS(path, ca, crt, key string) (tlsConfig *tls.Config) {
	// 1. 加载 CA 根证书（用于验证 EMQX 服务端身份）
	caCert, err := os.ReadFile(path + ca)
	if err != nil {
		log.Fatalf("读取 CA 证书失败: %v", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("解析 CA 证书失败")
	}
	// 2. 加载客户端证书和私钥（用于向 EMQX 证明自己的身份）
	clientCert, err := tls.LoadX509KeyPair(path+crt, path+key)
	if err != nil {
		log.Fatalf("加载客户端证书/私钥失败: %v", err)
	}

	// 3. 配置 TLS 参数
	tlsConfig = &tls.Config{
		RootCAs:      caCertPool,                    // 信任的 CA
		Certificates: []tls.Certificate{clientCert}, // 客户端证书
		// 注意：生产环境必须保持为 false（默认值），以验证服务端证书
		// 仅在开发测试自签名证书且遇到主机名不匹配时，可临时设为 true
		InsecureSkipVerify: false,
	}
	return tlsConfig
}
