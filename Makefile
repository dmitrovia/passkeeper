generate-keys-rsa:
	go run generate/generate.go;
	mv generate/public.pem internal/client/crypto/keys;
	mv generate/private.pem internal/server/crypto/keys;
	go run generate/generate.go;
	mv generate/private.pem internal/client/crypto/keys;
	mv generate/public.pem internal/server/crypto/keys;
build-client:
	go build -o ./cmd/client/main ./cmd/client/.
	mv main cmd/client
build-server:
	go build -o ./cmd/server/main ./cmd/server/.
	mv main cmd/server
generate-cert:
	openssl req -new -newkey rsa:2048 -keyout ca.key -x509 -sha256 -days 365 -out ca.crt -subj "/C=US/ST=California/L=Mountain View/O=Your Organization/OU=Your Unit/CN=localhost";
	openssl genrsa -out server.key 2048;
	openssl req -new -key server.key -out server.csr -config server.cnf;
	openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
		-CAcreateserial -out server.crt -days 365 -sha256 -extfile server.cnf -extensions v3_ext;
	mv server.crt internal/server/tls;
	mv server.key internal/server/tls;
	mv ca.crt internal/client/tls;