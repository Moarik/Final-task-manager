.PHONY: run
run:
	docker compose up --build

proto-user:
	protoc --proto_path=/usr/local/include --proto_path=proto --go_out=proto/gen --go-grpc_out=proto/gen proto/user/user_service.proto
proto-task:
	protoc --proto_path=/usr/local/include --proto_path=proto --go_out=proto/gen --go-grpc_out=proto/gen proto/task/task_service.proto
proto-statistics:
	protoc --proto_path=/usr/local/include --proto_path=proto --go_out=proto/gen --go-grpc_out=proto/gen proto/statistics/statistics_service.proto