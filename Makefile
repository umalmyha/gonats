tickets-gen:
	protoc --proto_path=./tickets-service --go_out=./tickets-service/rpc/ticket --twirp_out=./tickets-service/rpc/ticket tickets-service/rpc/ticket/service.proto
