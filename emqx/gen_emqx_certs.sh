#!/bin/bash

# 设置变量
CA_DIR="ca"
SERVER_DIR="server"
CLIENT_DIR="client"
DAYS=3650  # 证书有效期（10年）

echo "🚀 开始生成 EMQX mTLS 所需证书..."

# 1. 创建目录
mkdir -p $CA_DIR $SERVER_DIR $CLIENT_DIR

# 2. 生成 CA 根证书
echo "🔐 正在生成 CA 根证书..."
openssl genrsa -out $CA_DIR/ca.key 2048
openssl req -x509 -new -nodes -key $CA_DIR/ca.key -sha256 -days $DAYS \
  -subj "/C=CN/ST=Guangdong/L=Shenzhen/O=EMQX/CN=EMQX Root CA" \
  -out $CA_DIR/ca.crt

# 3. 生成 EMQX 服务端证书
echo "🖥️  正在生成 EMQX 服务端证书..."
openssl genrsa -out $SERVER_DIR/server.key 2048
openssl req -new -key $SERVER_DIR/server.key \
  -subj "/C=CN/ST=Guangdong/L=Shenzhen/O=EMQX/CN=EMQX Server" \
  -out $SERVER_DIR/server.csr

# 为服务端证书添加 SAN (Subject Alternative Name) 扩展
cat <<EOF > $SERVER_DIR/server_ext.cnf
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = emqx
IP.1 = 127.0.0.1
IP.2 = 172.17.0.2
EOF

openssl x509 -req -in $SERVER_DIR/server.csr -CA $CA_DIR/ca.crt -CAkey $CA_DIR/ca.key \
  -CAcreateserial -out $SERVER_DIR/server.crt -days $DAYS -sha256 \
  -extfile $SERVER_DIR/server_ext.cnf

# 4. 生成客户端证书 (用于测试或设备连接)
echo "📱 正在生成客户端证书..."
openssl genrsa -out $CLIENT_DIR/client.key 2048
openssl req -new -key $CLIENT_DIR/client.key \
  -subj "/C=CN/ST=Guangdong/L=Shenzhen/O=EMQX/CN=EMQX Client" \
  -out $CLIENT_DIR/client.csr
openssl x509 -req -in $CLIENT_DIR/client.csr -CA $CA_DIR/ca.crt -CAkey $CA_DIR/ca.key \
  -CAcreateserial -out $CLIENT_DIR/client.crt -days $DAYS -sha256

# 5. 清理临时文件
rm -f $SERVER_DIR/server.csr $CLIENT_DIR/client.csr $SERVER_DIR/server_ext.cnf $CA_DIR/ca.srl

echo "✅ 证书生成完毕！文件结构如下："
echo "📁 ca/        -> CA 根证书 (ca.crt, ca.key)"
echo "📁 server/    -> EMQX 服务端证书 (server.crt, server.key)"
echo "📁 client/    -> 客户端证书 (client.crt, client.key)"