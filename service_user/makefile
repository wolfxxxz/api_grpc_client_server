server_run:
	go run cmd/main.go
client_run:
	go run cmd/client/main.go
mock_user:
	~/go/bin/mockgen -source=internal/usecase/repository/user_repository.go -destination=./internal/mock/usecase_user_repository.go -package=mock
mock_cache:
	~/go/bin/mockgen -source=internal/usecase/cache/user_cache.go -destination=./internal/mock/usecase_user_cache.go -package=mock
run_bases:
	docker-compose up --build
get_grpc:
	go get -u google.golang.org/grpc
get_protobuf:
	google.golang.org/protobuf
proto_generate:
	protoc -I api/proto --go_out=plugins=grpc:internal/service_user api/proto/user_service.proto
