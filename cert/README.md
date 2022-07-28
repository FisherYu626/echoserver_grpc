###普通的CA证书

openssl genrsa -des3  -out ca.key 2048

openssl req -new -key ca.key -out ca.csr

openssl x509 -req -days 365 -in ca.csr -signkey ca.key -out ca.crt


###自go 1.17以后
需要使用san证书
san证书需要配置相应的自定义域名


openssl genpkey -algorithm RSA -out server.key

openssl req -new -nodes -key server.key -out server.csr -days 3650 -config ./openssl.cnf -extensions v3_req

openssl x509 -req -days 365 -in server.csr -out server.pem -CA ca.crt -CAkey ca.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req

证书后缀 windows一般用 crt linux 用pem


###双向认证生成客户端证书
openssl genpkey -algorithm RSA -out client.key
openssl req -new -nodes -key client.key -out client.csr -days 3650 -config ./openssl.cnf -extensions v3_req
openssl x509 -req -days 3650 -in client.csr -out client.pem -CA ca.crt -CAkey ca.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req