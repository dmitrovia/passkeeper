generate-key:
	go run generate/generate.go
move-key1:
	mv generate/public.pem internal/client/crypto/keys
	mv generate/private.pem internal/server/crypto/keys
move-key2:
	mv generate/private.pem internal/client/crypto/keys
	mv generate/public.pem internal/server/crypto/keys
build-client:
	go build -o ./cmd/client/main ./cmd/client/.
	mv main cmd/client
build-server:
	go build -o ./cmd/server/main ./cmd/server/.
	mv main cmd/server